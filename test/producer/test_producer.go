package main

import (
	"context"
	"flag"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// 向kafka topic异步写多条消息
// go test .\test\consumer_test.go -v -run TestProducerAsync -args -topic my-topic -address 127.0.0.1:9092 -user user -password pass
func main() { // TestProducerAsync(t *testing.T) {
	var topic = flag.String("topic", "my-topic", "kafka topic")
	var address = flag.String("address", "127.0.0.1:9092", "")
	var user = flag.String("user", "", "")
	var password = flag.String("password", "", "")

	flag.Parse()
	log.Println("topic: ", *topic)
	log.Println("address: ", *address)
	log.Println("user: ", *user)

	mechanism := plain.Mechanism{
		Username: *user,
		Password: *password,
	}
	// Transports are responsible for managing connection pools and other resources,
	// it's generally best to create a few of these and share them across your
	// application.
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	// 创建一个writer 向topic-A发送消息
	w := &kafka.Writer{
		Addr:         kafka.TCP(*address),
		Topic:        *topic,
		Transport:    sharedTransport,
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        true,                // 异步
	}

	// 一次写多条消息
	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
