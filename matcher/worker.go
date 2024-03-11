package matcher

import (
	"encoding/json"
	"sync"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/wasilibs/go-re2"

	"data-matcher/engine"
)

type Worker struct {
	Name            string               `json:"name"`
	Flowch          chan *nats.Msg       `json:"-"`
	Outch           chan *nats.Msg       `json:"-"`
	Dnsch           chan *DnsItem        `json:"-"`
	ValueRegs       []*engine.ValueRegex `json:"-"`
	ColDicts        []*engine.ColDict    `json:"-"`
	Appmap          *sync.Map            `json:"-"`
	Apimap          *sync.Map            `json:"-"`
	Ipmap           *sync.Map            `json:"-"`
	CountMsg        uint64               `json:"count_msg"`
	CountMatchRegex uint64               `json:"count_matched_regex"`
	CountMatchDict  uint64               `json:"count_matched_dict"`

	rs   []*re2.Regexp
	trie *ahocorasick.Trie
}

func (p *Worker) Init() error {
	p.rs = make([]*re2.Regexp, 0)
	for _, vreg := range p.ValueRegs {
		r, err := re2.Compile(vreg.VReg)
		if err != nil {
			return err
		}
		p.rs = append(p.rs, r)
	}

	dicts := make([]string, 0)
	for _, cdict := range p.ColDicts {
		dicts = append(dicts, cdict.CDict)
	}
	p.trie = ahocorasick.NewTrieBuilder().AddStrings(dicts).Build()
	// log.Debugf("worker [%s] init value regexs %d, column dicts %d", p.Name, len(p.ValueRegs), len(p.ColDicts))

	return nil
}

func (p *Worker) Run() {
	log.Infof("worker running with name [%s], value regexs %d, column dicts %d",
		p.Name, len(p.ValueRegs), len(p.ColDicts))
	p.CountMsg, p.CountMatchRegex, p.CountMatchDict = 0, 0, 0

	httpch := make(chan interface{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.processHttp(httpch)
	}()

	for m := range p.Flowch {
		p.CountMsg++
		// log.Debugf(m.Size(), len(m.Data))

		// 处理消息
		switch m.Subject {
		case "flow.http":
			// if matched := p.proccessHttp(m); matched {
			// 	p.Outch <- m
			// }
			httpch <- m
		case "flow.dns":
			if dnsitem, err := p.processDns(m); err == nil {
				p.Dnsch <- dnsitem
			}
		default:
		}

		if p.CountMsg%1000 == 0 {
			log.Debugf("worker [%s] msgs: %d, value regex matched %d, column dict matched %d",
				p.Name, p.CountMsg, p.CountMatchRegex, p.CountMatchDict)
		}
	}

	close(httpch)
	wg.Wait()
}

func (p *Worker) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
