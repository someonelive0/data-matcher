package matcher

import (
	"data-matcher/utils"
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type Outputer struct {
	Outch       chan *nats.Msg              `json:"-"`
	Dnsch       chan map[string]interface{} `json:"-"`
	NatsConfig  *NatsConfig                 `json:"-"`
	Stats       *MyStatistic                `json:"-"`
	CountMsg    uint64                      `json:"count_msg"`
	CountFailed uint64                      `json:"count_failed"`

	nc  *nats.Conn
	js  nats.JetStreamContext
	kvb nats.KeyValue
}

func (p *Outputer) init() error {
	servers := strings.Join(p.NatsConfig.Servers, ",")
	nc, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password)
	if err != nil {
		log.Errorf("ouputer NatsConnect %s failed: %s", servers, err)
		return err
	}

	// 创建JetStream上下文
	js, err := nc.JetStream(
		nats.PublishAsyncMaxPending(256),
		nats.PublishAsyncErrHandler(func(_ nats.JetStream, _ *nats.Msg, err error) { // 异步发布消息错误
			// TODO, 应该保存发布失败的消息，好下次发送
			log.Errorf("nats jetstream ErrorHandler error: %v", err)
			p.CountFailed++
		}),
	)
	if err != nil {
		log.Errorf("ouputer new jetstream %s failed: %s", servers, err)
		return err
	}
	// 创建流，流和主题都以"match_flow"开头，以区分kafka来的原始流
	if err = utils.NatsCreateStream(js, "match_flow", "match_flow.*"); err != nil {
		log.Errorf("ouputer create nats stream failed %s", err)
		return err
	}

	kvb, err := js.KeyValue("dns")
	if err != nil {
		kvb, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "dns",
		})
		if err != nil {
			log.Errorf("ouputer create nats keyvalue failed %s", err)
			return err
		}
	}

	p.nc = nc
	p.js = js
	p.kvb = kvb
	log.Infof("ouputer connect %s success by user %s", servers, p.NatsConfig.User)

	return nil
}

func (p *Outputer) Run() error {
	err := p.init()
	if err != nil {
		return err
	}

	httpch_closed, dnsch_closed := false, false
	for {
		if httpch_closed && dnsch_closed {
			return nil
		}

		select {
		// 输出已经匹配的http流
		case m, ok := <-p.Outch:
			if !ok {
				httpch_closed = true
			} else {
				_, err = p.js.PublishAsync("match_"+m.Subject, m.Data) // 异步发布
				if err != nil {
					p.CountFailed++
					log.Errorf("ouputer jetstream pub failed: %s", err)
				} else {
					p.Stats.OutputCount(1)
				}
			}

		// 输出DNS键值对
		case dnsmap, ok := <-p.Dnsch:
			if !ok {
				dnsch_closed = true
			} else {
				if key, ok := dnsmap["rrname"]; ok {
					b, _ := json.Marshal(dnsmap)
					if _, err = p.kvb.Put(key.(string), b); err != nil {
						log.Errorf("ouputer set kv failed: %s", err)
					}
				}
			}
		}
	}

	return nil
}

func (p *Outputer) Stop() error {
	if p.js != nil {
		<-p.js.PublishAsyncComplete() // should wait async publish finished
	}
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
