package matcher

import (
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

type Inputer struct {
	nc  *nats.Conn
	sub *nats.Subscription
}

// async to receive messages
func (p *Inputer) Run(msgch chan *nats.Msg,
	arg_server, arg_user, arg_password, arg_subject, arg_queue string) error {

	nc, err := utils.NatsConnect(arg_server, arg_user, arg_password)
	if err != nil {
		log.Errorf("Inputer NatsConnect failed: %s", err)
		return err
	}
	log.Infof("connect %s success by user %s\n", arg_server, arg_user)

	sub, err := utils.QueueSub2Chan(nc, arg_subject, arg_queue, msgch)
	if err != nil {
		log.Errorf("Inputer QueueSub2Chan failed: %s", err)
		nc.Close()
		return err
	}
	log.Infof("Inputer sub %s success with queue %s\n", arg_subject, arg_queue)

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
