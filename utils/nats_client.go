package utils

import (
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
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
		nats.Name("data-matcher 3.0.0"),
		nats.UserInfo(user, password),
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
		// nats.NoEcho(),
		nats.ReconnectWait(10*time.Second),
		nats.ReconnectBufSize(50*1024*1024),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			// handle disconnect error event
			if err != nil {
				log.Errorf("nats DisconnectErrHandler client disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// handle reconnect event
			log.Infof("nats ReconnectHandler client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Debugf("nats ClosedHandler client closed")
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Warnf("nats DiscoveredServersHandler client discover")
			log.Infof("nats Known servers: %v\n", nc.Servers())
			log.Infof("nats Discovered servers: %v\n", nc.DiscoveredServers())
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Errorf("nats ErrorHandler error: %v", err)
			// if err == nats.ErrSlowConsumer { // logSlowConsumer
			// pendingMsgs, _, err := sub.Pending()
			// if err != nil {
			// 	fmt.Printf("couldn't get pending messages: %v", err)
			// 	return
			// }
			// fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			// 	pendingMsgs, sub.Subject)
			// // Log error, notify operations...
			// }
		}),
	)
	if err != nil {
		// log.Errorf("NatsConnect failed: %s", err)
		return nil, err
	}

	// Do something with the connection
	mp := nc.MaxPayload()
	log.Debugf("nats maximum payload is %v bytes", mp)

	getStatusTxt := func(nc *nats.Conn) string {
		switch nc.Status() {
		case nats.CONNECTED:
			return "Connected"
		case nats.CLOSED:
			return "Closed"
		default:
			return "Other"
		}
	}
	log.Debugf("nats connection status is %v", getStatusTxt(nc))

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

func NatsSendMsg(nc *nats.Conn, subject, msg string) error {
	if err := nc.Publish(subject, []byte(msg)); err != nil {
		return err
	}

	return nil
}

// 查看流信息，如果不存在就创建
func NatsCreateStream(js nats.JetStreamContext, stream_name, stream_subjects string) error {
	stream, err := js.StreamInfo(stream_name)

	// stream not found, create it
	if err != nil || stream == nil {
		log.Infof("nats creating stream that is not existed, %s", stream_name)

		_, err = js.AddStream(&nats.StreamConfig{
			Name:              stream_name,
			Subjects:          []string{stream_subjects},
			Storage:           nats.FileStorage,
			Replicas:          1,
			Retention:         nats.LimitsPolicy,
			Discard:           nats.DiscardOld,
			MaxAge:            time.Duration(30 * 24 * time.Hour), // 30d
			MaxBytes:          10 * 1024 * 1024 * 1024,            // 10GiB
			MaxMsgSize:        1024 * 1024,                        // 1MiB
			MaxMsgs:           -1,
			MaxMsgsPerSubject: -1,
			Duplicates:        time.Duration(120 * time.Second),
			AllowRollup:       false,
			AllowDirect:       true,
			DenyDelete:        false,
			DenyPurge:         false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
