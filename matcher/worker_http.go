package matcher

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

// process msg of subject flow.http
func (p *Worker) processHttp(m *nats.Msg) bool {
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
