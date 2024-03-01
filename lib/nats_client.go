package mq

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// func opts() {
// 	servers := []string{"nats://" + *arg_host + ":4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}
// 	opts := nats.GetDefaultOptions()
// 	opts.Url = strings.Join(servers, ",")
// 	opts.Verbose = true
// 	opts.Pedantic = true
// }

// params:
// servers: servers := strings.Join([]string{"nats://127.0.0.1:4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}, ",")
//
//	or "nats://127.0.0.1:4222,nats://127.0.0.1:1223,nats://127.0.0.1:1224"
//
// user: nats user
// user: nats password of user
func NatsConnect(servers, user, password string) (*nats.Conn, error) {

	nc, err := nats.Connect(
		servers,
		nats.Name("natscli of mime 0.1.0"),
		nats.UserInfo(user, password),
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
		nats.NoEcho(),
		nats.ReconnectWait(10*time.Second),
		nats.ReconnectBufSize(5*1024*1024),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			// handle disconnect error event
			log.Printf("DisconnectErrHandler client disconnected: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// handle reconnect event
			log.Println("ReconnectHandler client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("ClosedHandler client closed")
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Printf("DiscoveredServersHandler client discover")
			log.Printf("Known servers: %v\n", nc.Servers())
			log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("ErrorHandler Error: %v", err)
		}), // logSlowConsumer
	)
	if err != nil {
		log.Fatal("Connect failed: ", err)
		return nil, err
	}

	return nc, nil
}

// Create a queue subscription on "updates" with queue name "workers"
// async subscription in a coroutine.
func QueueSub2Chan(nc *nats.Conn, subject, queue_name string, ch chan *nats.Msg) (*nats.Subscription, error) {
	// ch := make(chan *nats.Msg)
	sub, err := nc.QueueSubscribe(subject, queue_name, func(m *nats.Msg) {
		ch <- m
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}
