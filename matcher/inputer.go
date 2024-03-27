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
	CountSlow  uint64         `json:"count_slow"`

	nc  *nats.Conn
	sub *nats.Subscription
}

// async to receive messages from a flow
// 采用nats的subject通配符，可以在一个inputer中输入多个flow主题，即 subject: flow.*
func (p *Inputer) Run() error {
	natsErrHandler := func(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
		if natsErr == nats.ErrSlowConsumer { // 当出现慢消费者错误，即当前进程处理不了nats的消息，就计入统计，为以后扩展缓存和节点提供依据
			p.CountSlow++
			p.Stats.InputSlowCount(1)
			if p.CountSlow == 1 { // 只打印第一次慢消费者错误，避免错误日志过多
				log.Errorf("nats ErrorHandler error: %v", natsErr)
			}
		} else {
			log.Errorf("nats ErrorHandler error: %v", natsErr)
		}
	}

	servers := strings.Join(p.NatsConfig.Servers, ",")
	nc, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password, natsErrHandler)
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
