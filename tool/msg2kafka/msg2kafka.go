package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var (
	arg_broker   = flag.String("broker", "localhost:9092", "kafka server hostname or ip")
	arg_user     = flag.String("user", "", "kafka user")
	arg_password = flag.String("password", "", "")
	arg_topic    = flag.String("topic", "foo", "topic of kafka to receive message")
	arg_num      = flag.Int("num", 10, "message numbers, default 10 and at least 10")
	msg          = []byte(`{"timestamp":"2024-02-19T14:40:44.570856+0800","response_time":1,"flow_id":1633165473693756,"event_type":"http","pcap_filename":"","src_mac":"e8:65:49:4b:46:55","src_ip":"10.10.20.113","src_port":35616,"dst_mac":"00:27:0d:ff:5b:e1","dest_ip":"113.107.238.154","dest_port":8199,"proto":"TCP","tx_id":1,"http":{"hostname":"www.zs-hospital.sh.cn","http_port":8199,"url":"\/","http_user_agent":"Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)","http_method":"GET","protocol":"HTTP\/1.1","status":403,"length":2209,"request_headers":[{"name":"Host","value":"www.zs-hospital.sh.cn:8199"},{"name":"Accept-Charset","value":"iso-8859-1,utf-8;q=0.9,*;q=0.1"},{"name":"Accept-Language","value":"en"},{"name":"Connection","value":"Keep-Alive"},{"name":"User-Agent","value":"Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)"},{"name":"Pragma","value":"no-cache"},{"name":"Accept","value":"image\/gif, image\/x-xbitmap, image\/jpeg, image\/pjpeg, image\/png, *\/*"}],"response_headers":[{"name":"Date","value":"Fri, 18 Oct 2019 09:26:42 GMT"},{"name":"Connection","value":"keep-alive, close"},{"name":"Content-Length","value":"2209"},{"name":"X-Via-JSL","value":"76165a1,-"},{"name":"Set-Cookie","value":"__jsluid_h=dbab5e362df7eec9fba64290e2830ccf; max-age=31536000; path=\/; HttpOnly"},{"name":"X-Cache","value":"error"}],"reqHeaders":"Host: www.zs-hospital.sh.cn:8199\r\nAccept-Charset: iso-8859-1,utf-8;q=0.9,*;q=0.1\r\nAccept-Language: en\r\nConnection: Keep-Alive\r\nUser-Agent: Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)\r\nPragma: no-cache\r\nAccept: image\/gif, image\/x-xbitmap, image\/jpeg, image\/pjpeg, image\/png, *\/*\r\n\r\n","respHeaders":"Date: Fri, 18 Oct 2019 09:26:42 GMT\r\nConnection: keep-alive\r\nContent-Length: 2209\r\nConnection: close\r\nX-Via-JSL: 76165a1,-\r\nSet-Cookie: __jsluid_h=dbab5e362df7eec9fba64290e2830ccf; max-age=31536000; path=\/; HttpOnly\r\nX-Cache: error\r\n\r\n","reqBody":"","reqLen":0,"respBody":"","respLen":0,"totalLen":536},"host":"172.16.0.230"}`)
)

func init() {
	flag.Parse()
}

func main() {
	mechanism := plain.Mechanism{
		Username: *arg_user,
		Password: *arg_password,
	}
	// Transports are responsible for managing connection pools and other resources,
	// it's generally best to create a few of these and share them across your
	// application.
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	// 创建一个writer 向topic-A发送消息
	w := &kafka.Writer{
		Addr:         kafka.TCP(*arg_broker),
		Topic:        *arg_topic,
		Transport:    sharedTransport,
		Balancer:     &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks: kafka.RequireAll,    // ack模式
		Async:        true,                // 异步
	}

	log.Printf("sending %d messages with length %d\n", *arg_num, len(msg))
	t0 := time.Now()
	for i := 0; i < *arg_num/10; i++ {
		// 一次写十条消息
		err := w.WriteMessages(context.Background(),
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
			kafka.Message{Value: msg},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Printf("end sent %d messages in %v\n", *arg_num, time.Now().Sub(t0))
}
