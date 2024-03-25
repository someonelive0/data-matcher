package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"data-matcher/model"
	"data-matcher/utils"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type PostWorker struct {
	Httpch     chan *model.FlowHttp `json:"-"`
	NatsConfig *NatsConfig          `json:"-"`
	AppMap     sync.Map             `json:"-"` // app, number
	ApiMap     sync.Map             `json:"-"` // api, number
	IpMap      sync.Map             `json:"-"` // client ip, number

	CountApp    uint64 `json:"count_app"`
	CountApi    uint64 `json:"count_api"`
	CountIp     uint64 `json:"count_ip"`
	CountFailed uint64 `json:"count_failed"`
	nckv        *nats.Conn
	jskv        nats.JetStreamContext
	appkvb      nats.KeyValue
	ipkvb       nats.KeyValue
}

func (p *PostWorker) init() error {
	servers := strings.Join(p.NatsConfig.Servers, ",")

	// 写app,api,ip到key value store
	nckv, err := utils.NatsConnect(servers, p.NatsConfig.User, p.NatsConfig.Password)
	if err != nil {
		log.Errorf("post_worker NatsConnect kv %s failed: %s", servers, err)
		return err
	}
	jskv, err := nckv.JetStream( // 为KV创建JetStream上下文
		nats.PublishAsyncMaxPending(256),
		nats.PublishAsyncErrHandler(func(_ nats.JetStream, _ *nats.Msg, err error) { // 异步发布消息错误
			// TODO, 应该保存发布失败的消息，好下次发送
			log.Errorf("nats jetstream kv ErrorHandler error: %v", err)
			p.CountFailed++
		}),
	)
	if err != nil {
		log.Errorf("post_worker new jetstream kv %s failed: %s", servers, err)
		return err
	}
	// 创建名为app,api,ip的 KeyValue Bucket
	appkvb, err := jskv.KeyValue("app")
	if err != nil {
		appkvb, err = jskv.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:       "app",
			Storage:      nats.FileStorage,
			Replicas:     1,
			MaxBytes:     -1, // 1 * 1024 * 1024 * 1024, // 1GiB
			MaxValueSize: -1, // 1024 * 1024,            // 1MiB,
		})
		if err != nil {
			log.Errorf("post_worker create nats keyvalue failed %s", err)
			return err
		}
	}

	ipkvb, err := jskv.KeyValue("ip")
	if err != nil {
		ipkvb, err = jskv.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:       "ip",
			Storage:      nats.FileStorage,
			Replicas:     1,
			MaxBytes:     -1, // 1 * 1024 * 1024 * 1024, // 1GiB
			MaxValueSize: -1, // 1024 * 1024,            // 1MiB,
		})
		if err != nil {
			log.Errorf("post_worker create nats keyvalue failed %s", err)
			return err
		}
	}
	p.nckv = nckv
	p.jskv = jskv
	p.appkvb = appkvb
	p.ipkvb = ipkvb

	log.Infof("post_worker connect %s success by user %s", servers, p.NatsConfig.User)

	return nil
}

func (p *PostWorker) Run() error {
	err := p.init()
	if err != nil {
		return err
	}

	// Fan out msg to all msg processer
	var wg sync.WaitGroup

	// 发现APP协程
	var discoverAppChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.discoverApp(flowHttp)
		}
	}(discoverAppChan)

	// 发现API协程
	var discoverApiChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.discoverApi(flowHttp)
		}
	}(discoverApiChan)

	// 发现IP协程
	var discoverIpChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.discoverIp(flowHttp)
		}
	}(discoverIpChan)

	// 发现Account协程
	var discoverAccountChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.discoverAccount(flowHttp)
		}
	}(discoverAccountChan)

	// 获取Token协程
	var disposeTokenChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.disposeToken(flowHttp)
		}
	}(disposeTokenChan)

	// 获取Username协程
	var disposeUsernameChan = make(chan *model.FlowHttp)
	wg.Add(1)
	go func(ch chan *model.FlowHttp) {
		defer wg.Done()
		for flowHttp := range ch {
			p.disposeUsername(flowHttp)
		}
	}(disposeUsernameChan)

	for flowHttp := range p.Httpch {
		discoverAppChan <- flowHttp
		discoverApiChan <- flowHttp
		discoverIpChan <- flowHttp
		discoverAccountChan <- flowHttp
		disposeTokenChan <- flowHttp
		disposeUsernameChan <- flowHttp
	}

	// 当p.Httpch被关闭后，关闭下游的所有协程channel，使协程退出
	close(discoverAppChan)
	close(discoverApiChan)
	close(discoverIpChan)
	close(discoverAccountChan)
	close(disposeTokenChan)
	close(disposeUsernameChan)
	wg.Wait()

	return nil
}

func (p *PostWorker) Stop() {
	if p.jskv != nil {
		<-p.jskv.PublishAsyncComplete() // should wait async publish finished
	}
	if p.nckv != nil {
		p.nckv.Close()
		p.nckv = nil
	}
}

// discover app from msg of subject flow.http
func (p *PostWorker) discoverApp(flowHttp *model.FlowHttp) error {

	// 由于消息中的hostname可能含有路径，需求去除，例如  "hostname": "mzzj.sh.gov.cn/..",
	// 或者 www.oldkids.cn?<script>cross_site_scripting.nasl
	// 或者 hostname='
	hostname := flowHttp.Http.Hostname
	if n := strings.IndexAny(hostname, "/?='<%("); n != -1 {
		hostname = hostname[:n]
	}

	app := fmt.Sprintf("%s/%s/%d", hostname, flowHttp.Dest_ip, flowHttp.Dest_port)
	_, ok := p.AppMap.Load(app)
	if !ok {
		p.AppMap.Store(app, 1)
		p.CountApp++

		if false {
			// 记录本地Map同时写入nats keyvalue store
			if _, err := p.appkvb.Put(app, []byte{}); err != nil {
				// p.CountDnsFailed++
				log.Errorf("post_worker set app kv [%s] failed: %s", app, err)
			}
		}
	}

	return nil
}

// discover api from msg of subject flow.http
func (p *PostWorker) discoverApi(flowHttp *model.FlowHttp) error {
	api := flowHttp.Http.Url
	i := strings.IndexByte(api, '?') // erase ?
	if i != -1 {
		api = api[:i]
	}
	i = strings.IndexByte(api, '#') // erase #
	if i != -1 {
		api = api[:i]
	}

	api = flowHttp.Http.Hostname + api
	_, ok := p.ApiMap.Load(api)
	if !ok {
		p.ApiMap.Store(api, 1)
		p.CountApi++
	}
	return nil
}

// dispose username from msg of subject flow.http
func (p *PostWorker) disposeUsername(flowHttp *model.FlowHttp) error {
	return nil
}

// dispose token from msg of subject flow.http
func (p *PostWorker) disposeToken(flowHttp *model.FlowHttp) error {
	return nil
}

// discover account from msg of subject flow.http
func (p *PostWorker) discoverAccount(flowHttp *model.FlowHttp) error {
	return nil
}

// discover client ip from msg of subject flow.http
func (p *PostWorker) discoverIp(flowHttp *model.FlowHttp) error {
	if len(flowHttp.Src_ip) == 0 {
		return nil
	}

	_, ok := p.IpMap.Load(flowHttp.Src_ip)
	if !ok {
		p.IpMap.Store(flowHttp.Src_ip, 1)
		p.CountIp++

		if false {
			// 记录本地Map同时写入nats keyvalue store
			if _, err := p.ipkvb.Put(flowHttp.Src_ip, []byte{}); err != nil {
				// p.CountDnsFailed++
				log.Errorf("post_worker set ip kv [%s] failed: %s", flowHttp.Src_ip, err)
			}
		}
	}

	return nil
}

func (p *PostWorker) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
