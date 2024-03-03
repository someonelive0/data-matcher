package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

// 从单个主题-分区（topic-partition）消费消息这种典型场景
func TestConsumer(t *testing.T) {
	// 指定要连接的topic和partition
	topic := "`my-topic`"
	partition := 0

	// 连接至Kafka的leader节点
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// 读取一批消息，得到的batch是一系列消息的迭代器
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	// 遍历读取消息
	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

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
func TestConsumerGroup(t *testing.T) {
	// 创建一个reader，指定GroupID，从 topic-A 消费消息
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		GroupID:  "consumer-group-id", // 指定消费者组id
		Topic:    "topic-A",
		MaxBytes: 10e6, // 10MB
	})

	ctx := context.Background()
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
