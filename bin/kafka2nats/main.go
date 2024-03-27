package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

var (
	arg_debug   = flag.Bool("D", false, "debug")
	arg_version = flag.Bool("v", false, "version")
	arg_config  = flag.String("f", "etc/kafka2nats.yaml", "config filename")
	START_TIME  = time.Now()
)

func init() {
	flag.Parse()
	if *arg_version {
		fmt.Print(utils.IDSS_BANNER)
		fmt.Printf("%s\n", utils.APP_VERSION)
		os.Exit(0)
	}
	utils.ShowBannerForApp("kafka2nats", utils.APP_VERSION, utils.BUILD_TIME)
	utils.Chdir2PrgPath()
	fmt.Println("pwd:", utils.GetPrgDir())
	utils.InitLog("kafka2nats.log", *arg_debug)
	log.Infof("BEGIN... %v, config=%v, debug=%v",
		START_TIME.Format(time.DateTime), *arg_config, *arg_debug)
}

func main() {
	myconfig, err := LoadConfig(*arg_config)
	if err != nil {
		log.Errorf("read config failed: %v", err)
		os.Exit(1)
	}
	log.Infof("myconfig: %s", myconfig.Dump())

	var ch = make(chan []byte)
	var wg sync.WaitGroup
	var ctx, cancel = context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		output(ch, &myconfig.NatsConfig, myconfig.Kafka2Nats.Subject)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		input(ctx, ch, &myconfig.KafkaConfig, myconfig.Kafka2Nats.Topic)
	}()

	// wait until graceful exit
	signalChan := utils.NewShutdownSignal()
	utils.WaitExit(signalChan, func() {
		log.Infof("receive exit signal, exit...")
		// your clean code
		cancel()
		close(ch)
		wg.Wait()
	})
}

// 从kafka读入数据，写入channel
func input(ctx context.Context, ch chan []byte, kafkaConfig *KafkaConfig, topic string) error {
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: nil,
	}

	switch kafkaConfig.Mechanism {
	case "PLAIN":
		dialer.SASLMechanism = plain.Mechanism{
			Username: kafkaConfig.User,
			Password: kafkaConfig.Password,
		}
	case "SCRAM":
		mechanism, err := scram.Mechanism(scram.SHA512, kafkaConfig.User, kafkaConfig.Password)
		if err != nil {
			panic(err)
		}
		dialer.SASLMechanism = mechanism
	}

	// 创建一个reader，指定GroupID，从 topic-A 消费消息
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaConfig.Brokers,
		GroupID:  kafkaConfig.Group, // 指定消费者组id
		Topic:    topic,
		Dialer:   dialer,
		MaxBytes: 10e6, // 10MB
	})
	log.Infof("connect kafka %v with user [%s] success", kafkaConfig.Brokers, kafkaConfig.User)

	// TODO, 利用context，从main() 来取消阻塞执行的 FetchMessage，并退出
	// ctx := context.Background()
	for {
		// 获取消息
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				log.Debugf("kafka get context canceled:%v", err)
				break
			}
			if errors.Is(err, io.EOF) { // 当reader.Close后，进入这个分支
				log.Debugf("kafka get eof")
				break
			}

			log.Errorf("FetchMessage error: %s", err)

			// TODO, 考虑这里需要进行Kafka重连，而不是break循环
			time.Sleep(3 * time.Second)
			continue
		}

		// 处理消息
		//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		log.Debugf("message at topic/partition/offset %v/%v/%v: %s = %d", m.Topic, m.Partition, m.Offset, string(m.Key), len(m.Value))
		ch <- m.Value

		// 显式提交, 注释可以连续测试
		// if err := r.CommitMessages(ctx, m); err != nil {
		// 	log.Fatal("failed to commit messages:", err)
		// }
	}
	// log.Debugf("input end")

	return nil
}

// 从channel读入数据，写入nats
func output(ch chan []byte, natsConfig *NatsConfig, subject string) error {
	nc, err := utils.NatsConnect(strings.Join(natsConfig.Servers, ","), natsConfig.User, natsConfig.Password, nil)
	if err != nil {
		log.Errorf("connect nats %v with user [%s] failed: %s", natsConfig.Servers, natsConfig.User, err)
		return err
	}
	defer nc.Close()
	log.Infof("connect nats %v with user %s success", natsConfig.Servers, natsConfig.User)

	count := 0
	t0 := time.Now()
	for msg := range ch {
		if err := nc.Publish(subject, msg); err != nil {
			return err
		}
		count++

		if count%10 == 0 {
			if err = nc.Flush(); err != nil {
				log.Errorf("flush msg failed: %s", err)
			}
		}
	}

	if err = nc.Flush(); err != nil {
		log.Errorf("flush msg failed: %s", err)
	}
	log.Infof("end sent %d messages in %v", count, time.Since(t0))

	return nil
}
