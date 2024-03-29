package matcher

import (
	"data-matcher/model"
	"data-matcher/utils"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type Outputer struct {
	LabelHttpch    chan *model.FlowHttp `json:"-"`
	LabelDnsch     chan *model.FlowDns  `json:"-"`
	NatsConfig     *NatsConfig          `json:"-"`
	Stats          *MyStatistic         `json:"-"`
	CountMsg       uint64               `json:"count_msg"`
	CountFailed    uint64               `json:"count_failed"`
	CountDnsMsg    uint64               `json:"count_dns_msg"`
	CountDnsFailed uint64               `json:"count_dns_failed"`

	nc   *nats.Conn // 写http到jetstream
	js   nats.JetStreamContext
	nckv *nats.Conn // 写dns到key value store
	jskv nats.JetStreamContext
	kvb  nats.KeyValue
}

func (p *Outputer) init() error {
	servers := strings.Join(p.NatsConfig.Servers, ",")

	{ // 写http到jetstream
		nc, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password, nil)
		if err != nil {
			log.Errorf("ouputer NatsConnect %s failed: %s", servers, err)
			return err
		}
		js, err := nc.JetStream( // 创建JetStream上下文
			nats.PublishAsyncMaxPending(25600), // 增加异步写等待数可以防止 too many async write 错误
			nats.MaxWait(10*time.Second),
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
		p.nc = nc
		p.js = js
	}

	{ // 写dns到key value store

		nckv, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password, nil)
		if err != nil {
			log.Errorf("ouputer NatsConnect kv %s failed: %s", servers, err)
			return err
		}
		jskv, err := nckv.JetStream( // 为KV创建JetStream上下文
			nats.PublishAsyncMaxPending(25600),
			nats.MaxWait(10*time.Second),
		)
		if err != nil {
			log.Errorf("ouputer new jetstream kv %s failed: %s", servers, err)
			return err
		}
		// 创建名为dns的 KeyValue Bucket
		kvb, err := jskv.KeyValue("dns")
		if err != nil {
			kvb, err = jskv.CreateKeyValue(&nats.KeyValueConfig{
				Bucket:       "dns",
				Storage:      nats.FileStorage,
				Replicas:     1,
				MaxBytes:     -1, // 1 * 1024 * 1024 * 1024, // 1GiB
				MaxValueSize: -1, // 1024 * 1024,            // 1MiB,
			})
			if err != nil {
				log.Errorf("ouputer create nats keyvalue dns failed %s", err)
				return err
			}
		}
		p.nckv = nckv
		p.jskv = jskv
		p.kvb = kvb
	}

	log.Infof("ouputer connect %s success by user %s", servers, p.NatsConfig.User)

	return nil
}

func (p *Outputer) Run() error {
	err := p.init()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.OutputHttp()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.OutputDns()
	}()

	wg.Wait()

	return nil
}

func (p *Outputer) OutputHttp() (err error) {
	for flowHttp := range p.LabelHttpch {
		p.CountMsg++
		b, _ := json.Marshal(flowHttp)
		_, err = p.js.PublishAsync("match_flow.http", b) // 异步发布
		if err != nil {
			log.Warnf("ouputer jetstream async pub failed: %s", err)

			// retry pub if failed async pub
			if _, err = p.js.Publish("match_flow.http", b); err != nil { // 同步重试
				p.CountFailed++
				log.Errorf("ouputer jetstream async and sync pub failed: %s", err)
			} else {
				p.Stats.OutputHttpCount(1)
			}
		} else {
			p.Stats.OutputHttpCount(1)
		}
	}

	return nil
}

func (p *Outputer) OutputDns() error {
	for flowDns := range p.LabelDnsch {
		p.CountDnsMsg++
		if false {
			go func(flowDns *model.FlowDns) { // 写dns到nats keyvalue store 太慢，所以用异步写
				if _, err := p.kvb.Get(flowDns.Dns.Rrname); err != nil { // 如果key不存在才Put
					b, _ := json.Marshal(flowDns.Dns)
					if _, err = p.kvb.Create(flowDns.Dns.Rrname, b); err != nil {
						p.CountDnsFailed++
						log.Errorf("ouputer set kv [%s] failed: %s", flowDns.Dns.Rrname, err)
					} else {
						p.Stats.OutputDnsCount(1)
					}
				}
			}(flowDns)
		}
	}

	return nil
}

func (p *Outputer) Stop() error {
	if p.js != nil {
		<-p.js.PublishAsyncComplete() // should wait async publish finished
	}
	if p.jskv != nil {
		<-p.jskv.PublishAsyncComplete() // should wait async publish finished
	}
	if p.nc != nil {
		p.nc.Close()
		p.nc = nil
	}
	if p.nckv != nil {
		p.nckv.Close()
		p.nckv = nil
	}
	return nil
}

func (p *Outputer) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
