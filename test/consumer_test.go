package test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var topic = flag.String("topic", "my-topic", "kafka topic")
var address = flag.String("address", "127.0.0.1:9092", "")
var user = flag.String("user", "", "")
var password = flag.String("password", "", "")

// 从单个主题-分区（topic-partition）消费消息这种典型场景
// go test .\test\consumer_test.go -v -run TestConsumer -args -topic my-topic -address 127.0.0.1:9092 -user user -password pass
func TestConsumer(t *testing.T) {
	flag.Parse()
	log.Println("topic: ", *topic)
	log.Println("address: ", *address)
	log.Println("user: ", *user)

	// 指定要连接的topic和partition
	partition := 0

	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	conn, err := dialer.DialLeader(context.Background(), "tcp", *address, *topic, partition)
	if err != nil {
		log.Fatal("failed to dial failed:", err)
	}
	log.Println("kafka connect success")

	// 连接至Kafka的leader节点
	// conn, err := dialer.DialLeader(context.Background(), "tcp", *address, topic, partition)
	// if err != nil {
	// 	log.Fatal("failed to dial leader:", err)
	// }

	// 设置读取超时时间，可以让测试程序退出循环
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	// 读取一批消息，得到的batch是一系列消息的迭代器
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	// 遍历读取消息
	count := 0
	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			log.Println("read kafka error: ", err)
			break
		}
		count++
		fmt.Println(string(b[:n]))
	}
	log.Println("read messages from kafka: ", count)

	// 关闭batch
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

// 消费整个Topic的所有消息这种典型场景，由broker管理的offset
// go test .\test\consumer_test.go -v -run TestConsumerGroup -args -topic my-topic -address 127.0.0.1:9092 -user user -password pass
func TestConsumerGroup(t *testing.T) {
	flag.Parse()
	log.Println("topic: ", *topic)
	log.Println("address: ", *address)
	log.Println("user: ", *user)

	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	// 创建一个reader，指定GroupID，从 topic-A 消费消息
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{*address},
		GroupID:  "consumer-test", // 指定消费者组id
		Topic:    *topic,
		Dialer:   dialer,
		MaxBytes: 10e6, // 10MB
	})

	// ctx := context.Background()
	// 使用 WithTimeout 创建一个带有超时的上下文对象
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 接收消息
	for {
		// m, err := r.ReadMessage(context.Background())
		// if err != nil {
		// 	break
		// }
		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		// 获取消息
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Println("FetchMessage error: ", err)
			break
		}
		// 处理消息
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 显式提交
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}

	// 程序退出前关闭Reader
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
