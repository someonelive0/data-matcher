package matcher

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

type Inputer struct {
	CountMsg uint64       `json:"count_msg"`
	Stats    *MyStatistic `json:"-"`

	nc  *nats.Conn
	sub *nats.Subscription
}

// async to receive messages
func (p *Inputer) Run(msgch chan *nats.Msg,
	arg_server, arg_user, arg_password, arg_subject, arg_queue string) error {

	nc, err := utils.NatsConnect(arg_server, arg_user, arg_password)
	if err != nil {
		log.Errorf("inputer NatsConnect failed: %s", err)
		return err
	}
	log.Infof("inputer connect %s success by user %s", arg_server, arg_user)

	// sub, err := utils.QueueSub2Chan(nc, arg_subject, arg_queue, msgch)
	sub, err := nc.QueueSubscribe(arg_subject, arg_queue, func(m *nats.Msg) {
		msgch <- m
		p.CountMsg++
		p.Stats.InputCount(1)
	})
	if err != nil {
		log.Errorf("inputer QueueSub2Chan failed: %s", err)
		nc.Close()
		return err
	}
	log.Infof("inputer sub %s success with queue %s", arg_subject, arg_queue)

	p.nc = nc
	p.sub = sub

	return nil
}

func (p *Inputer) Stop() error {
	if p.sub != nil {
		p.sub.Unsubscribe()
		p.sub = nil
	}
	if p.nc != nil {
		p.nc.Close()
		p.nc = nil
	}
	return nil
}

func (p *Inputer) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
