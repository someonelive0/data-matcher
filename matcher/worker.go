package matcher

import (
	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/wasilibs/go-re2"
)

type Worker struct {
	Name            string
	Msgch           chan *nats.Msg
	ValueRegs       []string
	ColDicts        []string
	CountMsg        int64
	CountMatchRegex int64
	CountMatchDict  int64

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

	for m := range p.Msgch {
		p.CountMsg++
		// log.Debugf(m.Size(), len(m.Data))

		// 依次匹配正则表达式
		for i, r := range p.rs {
			loc := r.FindIndex(m.Data)
			if len(loc) > 0 {
				log.Debugf("regex find rule %d with position %v", i, loc)
				p.CountMatchRegex++
			}
		}

		// 一次多模式匹配Dictionary
		matches := p.trie.Match(m.Data)
		p.CountMatchDict += int64(len(matches))

		if p.CountMsg%1000 == 0 {
			log.Infof("worker [%s] count: %d, matched value regex count %d, matched column dict count %d",
				p.Name, p.CountMsg, p.CountMatchRegex, p.CountMatchDict)
		}
	}
}
