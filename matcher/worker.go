package matcher

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	ahocorasick "github.com/BobuSumisu/aho-corasick"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/wasilibs/go-re2"

	"data-matcher/engine"
)

type Worker struct {
	Name            string                      `json:"name"`
	Flowch          chan *nats.Msg              `json:"-"`
	Outch           chan *nats.Msg              `json:"-"`
	Dnsch           chan map[string]interface{} `json:"-"`
	ValueRegs       []*engine.ValueRegex        `json:"-"`
	ColDicts        []*engine.ColDict           `json:"-"`
	CountMsg        uint64                      `json:"count_msg"`
	CountMatchRegex uint64                      `json:"count_matched_regex"`
	CountMatchDict  uint64                      `json:"count_matched_dict"`

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

	for m := range p.Flowch {
		p.CountMsg++
		// log.Debugf(m.Size(), len(m.Data))

		// 处理消息
		switch m.Subject {
		case "flow.http":
			if matched := p.proccessHttp(m); matched {
				p.Outch <- m
			}
		case "flow.dns":
			if dnsmap, err := p.proccessDns(m); err == nil {
				p.Dnsch <- dnsmap
			}
		default:
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

// process msg with subject flow.http
func (p *Worker) proccessHttp(m *nats.Msg) bool {
	// 只匹配 Response Data
	http, err := sonic.Get(m.Data, "http")
	if err != nil {
		// log.Errorf("worker unmarshal http failed: %s\n", err)
		return false
	}
	respBody, err := http.GetByPath("respBody").Raw()
	if err != nil || len(respBody) == 0 {
		return false
	}
	// log.Debugf("http resp %d: %s", len(respBody), respBody)
	respBodyb := []byte(respBody)

	// 依次匹配正则表达式
	reg_lable, reg_matched := p.matchReg(respBodyb)

	// 一次多模式匹配Dictionary
	dict_lable, dict_magched := p.matchDict(respBodyb)

	// 如果匹配到了正则或字典，写入输出队列
	if reg_matched || dict_magched {
		if len(reg_lable) > 0 {
			reg_lable = `, "value_regex": [` + reg_lable + `]`
			if m.Data[len(m.Data)-1] == '}' { // 删除最后一个大括号
				m.Data = append(m.Data[:len(m.Data)-1], []byte(reg_lable)...)
			} else {
				m.Data = append(m.Data, []byte(reg_lable)...)
			}
		}
		if len(dict_lable) > 0 {
			dict_lable = `, "column_dict": [` + dict_lable + `]`
			if m.Data[len(m.Data)-1] == '}' { // 删除最后一个大括号
				m.Data = append(m.Data[:len(m.Data)-1], []byte(dict_lable)...)
			} else {
				m.Data = append(m.Data, []byte(dict_lable)...)
			}
		}
		m.Data = append(m.Data, '}') // 添加最后一个大括号
		return true
	}

	return false
}

func (p *Worker) matchReg(data []byte) (string, bool) {

	// 依次匹配正则表达式
	reg_lable, reg_matched := "", false
	for i, r := range p.rs {
		loc := r.FindIndex([]byte(data))
		if loc != nil {
			// log.Debugf("regex find rule %d with position %v", i, loc)
			if len(reg_lable) > 0 {
				reg_lable += ", "
			}
			b, _ := json.Marshal(loc)
			reg_lable += fmt.Sprintf(`{"lable": "%s", "position": %s}`, p.ValueRegs[i].InfoType, b)
			atomic.AddUint64(&p.ValueRegs[i].CountMatch, 1)

			p.CountMatchRegex++
			reg_matched = true
		}
	}

	return reg_lable, reg_matched
}

func (p *Worker) matchDict(data []byte) (string, bool) {
	dict_lable, dict_magched := "", false
	matches := p.trie.Match(data)
	if len(matches) > 0 {
		for i, match := range matches {
			if i > 0 {
				dict_lable += ", "
			}
			dict_lable += fmt.Sprintf(`{"lable": "%s", "match": "%s", "position": %v}`,
				p.ColDicts[match.Pattern()].InfoType, match.Match(), match.Pos())
			atomic.AddUint64(&p.ColDicts[match.Pattern()].CountMatch, 1)
		}

		p.CountMatchDict += uint64(len(matches))
		dict_magched = true
	}

	return dict_lable, dict_magched
}
