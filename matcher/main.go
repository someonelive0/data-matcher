package main

import (
	"flag"
	"log"
	"strconv"
	"sync"

	mylib "data-matcher/lib"

	"github.com/nats-io/nats.go"
)

var arg_server = flag.String("server", "nats://localhost:4222", "nats server hostname or ip")
var arg_user = flag.String("user", "", "nats user")
var arg_password = flag.String("password", "", "")
var arg_subject = flag.String("subject", "foo", "subject of nats to receive message")
var arg_queue = flag.String("queue", "foo", "queue name for the subject")

func init() {
	flag.Parse()
}

func main() {
	nc, err := mylib.NatsConnect(*arg_server, *arg_user, *arg_password)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()
	log.Printf("connect %s success by user %s\n", *arg_server, *arg_user)

	msgch := make(chan *nats.Msg)
	sub, err := mylib.QueueSub2Chan(nc, *arg_subject, *arg_queue, msgch)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sub.Unsubscribe()
	log.Printf("sub %s success with queue %s\n", *arg_subject, *arg_queue)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			worker(name, msgch)
		}(strconv.Itoa(i))
	}

	wg.Wait()
}
