package matcher

import (
	"data-matcher/utils"
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type Outputer struct {
	Outch       chan *nats.Msg `json:"-"`
	NatsConfig  *NatsConfig    `json:"-"`
	HttpOutFlow *Flow          `json:"-"`
	Stats       *MyStatistic   `json:"-"`
	CountMsg    uint64         `json:"count_msg"`
	CountFailed uint64         `json:"count_failed"`

	nc *nats.Conn
}

func (p *Outputer) Run() error {
	servers := strings.Join(p.NatsConfig.Servers, ",")
	nc, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password)
	if err != nil {
		log.Errorf("ouputer NatsConnect %s failed: %s", servers, err)
		return err
	}

	// 创建JetStream上下文
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Errorf("ouputer new jetstream %s failed: %s", servers, err)
		return err
	}
	// 创建流，流和主题都以"match_flow"开头，以区分kafka来的原始流
	if err = utils.NatsCreateStream(js, "match_flow", "match_flow.*"); err != nil {
		log.Errorf("ouputer create nats stream failed %s", err)
		return err
	}
	p.nc = nc
	log.Infof("ouputer connect %s success by user %s", servers, p.NatsConfig.User)

	for m := range p.Outch {
		p.CountMsg++

		// _, err = js.Publish("match_"+m.Subject, m.Data) // 同步发布
		_, err = js.PublishAsync("match_"+m.Subject, m.Data) // 异步发布
		if err != nil {
			p.CountFailed++
			log.Errorf("ouputer jetstream pub failed: %s", err)
			continue
		}
		p.Stats.OutputCount(1)
	}

	return nil
}
func (p *Outputer) Stop() error {
	if p.nc != nil {
		p.nc.Close()
		p.nc = nil
	}
	return nil
}

func (p *Outputer) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
