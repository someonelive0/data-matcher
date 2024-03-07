package matcher

import (
	"encoding/json"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/wasilibs/go-re2"
)

type Worker struct {
	Name            string         `json:"name"`
	Msgch           chan *nats.Msg `json:"-"`
	Outch           chan *nats.Msg `json:"-"`
	ValueRegs       []string       `json:"-"`
	ColDicts        []string       `json:"-"`
	CountMsg        uint64         `json:"count_msg"`
	CountMatchRegex uint64         `json:"count_matched_regex"`
	CountMatchDict  uint64         `json:"count_matched_dict"`

	rs   []*re2.Regexp
	trie *ahocorasick.Trie
}

func (p *Worker) Init() error {
	p.rs = make([]*re2.Regexp, 0)
	for _, reg := range p.ValueRegs {
		r, err := re2.Compile(reg)
		if err != nil {
			return err
		}
		p.rs = append(p.rs, r)
	}

	p.trie = ahocorasick.NewTrieBuilder().AddStrings(p.ColDicts).Build()
	// log.Debugf("worker [%s] init value regexs %d, column dicts %d", p.Name, len(p.ValueRegs), len(p.ColDicts))

	return nil
}

func (p *Worker) Run() {
	log.Infof("worker running with name [%s], value regexs %d, column dicts %d",
		p.Name, len(p.ValueRegs), len(p.ColDicts))
	p.CountMsg = 0
	p.CountMatchRegex = 0
	p.CountMatchDict = 0

	var reg_matched, dict_magched = false, false
	for m := range p.Msgch {
		p.CountMsg++
		reg_matched, dict_magched = false, false
		// log.Debugf(m.Size(), len(m.Data))

		// 依次匹配正则表达式
		for _, r := range p.rs {
			loc := r.FindIndex(m.Data)
			if loc != nil {
				// log.Debugf("regex find rule %d with position %v", i, loc)
				p.CountMatchRegex++
				reg_matched = true
			}
		}

		// 一次多模式匹配Dictionary
		matches := p.trie.Match(m.Data)
		if len(matches) > 0 {
			p.CountMatchDict += uint64(len(matches))
			dict_magched = true
		}

		// 如果匹配到了正则或字典，写入输出队列
		if reg_matched || dict_magched {
			p.Outch <- m
		}

		if p.CountMsg%1000 == 0 {
			log.Debugf("worker [%s] msgs: %d, value regex matched %d, column dict matched %d",
				p.Name, p.CountMsg, p.CountMatchRegex, p.CountMatchDict)
		}
	}
}

func (p *Worker) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
