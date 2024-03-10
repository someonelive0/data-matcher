package matcher

import (
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

type Inputer struct {
	Flowch     chan *nats.Msg `json:"-"`
	NatsConfig *NatsConfig    `json:"-"`
	Flow       *Flow          `json:"-"`
	Stats      *MyStatistic   `json:"-"`
	CountMsg   uint64         `json:"count_msg"`

	nc  *nats.Conn
	sub *nats.Subscription
}

// async to receive messages from a flow
// 采用nats的subject通配符，可以在一个inputer中输入多个flow主题，即 subject: flow.*
func (p *Inputer) Run() error {
	servers := strings.Join(p.NatsConfig.Servers, ",")
	nc, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password)
	if err != nil {
		log.Errorf("inputer NatsConnect %s failed: %s", servers, err)
		return err
	}
	log.Infof("inputer connect %s success by user %s", servers, p.NatsConfig.User)

	sub, err := nc.QueueSubscribe(p.Flow.Subject, p.Flow.QueueName, func(m *nats.Msg) {
		p.Flowch <- m
		p.CountMsg++
		p.Stats.InputCount(1)
	})
	if err != nil {
		log.Errorf("inputer QueueSubscribe failed: %s", err)
		nc.Close()
		return err
	}
	log.Infof("inputer sub %s success with queue %s", p.Flow.Subject, p.Flow.QueueName)

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
