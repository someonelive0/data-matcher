package matcher

import (
	"sync/atomic"

	"data-matcher/model"
)

// process msg of subject flow.http
func (p *Worker) processMsgHttp(msgHttp *model.MsgHttp) bool {
	// 只匹配 Request Body 和 Response Body
	if msgHttp.Http.RespLen == 0 || len(msgHttp.Http.RespBody) == 0 {
		return false
	}
	respBodyb := []byte(msgHttp.Http.RespBody)

	// 依次匹配正则表达式
	msgHttp.Http.Value_regex = p.matchReg(respBodyb)

	// 一次多模式匹配Dictionary
	msgHttp.Http.Column_dict = p.matchDict(respBodyb)

	// 如果匹配到了正则或字典，写入输出队列
	if msgHttp.Http.Value_regex != nil || msgHttp.Http.Column_dict != nil {
		return true
	}

	return false
}

// 依次匹配正则表达式
func (p *Worker) matchReg(data []byte) []model.RegexMatched {
	var reg_matches []model.RegexMatched = nil

	for i, r := range p.rs {
		loc := r.FindIndex(data)
		if loc != nil {
			// log.Debugf("regex find rule %d with position %v", i, loc)
			matched := model.RegexMatched{
				Lable:    p.ValueRegs[i].InfoType,
				Position: loc,
			}
			if reg_matches == nil {
				reg_matches = make([]model.RegexMatched, 0)
			}
			reg_matches = append(reg_matches, matched)

			atomic.AddUint64(&p.ValueRegs[i].CountMatch, 1)
			p.CountMatchRegex++
		}
	}

	return reg_matches
}

func (p *Worker) matchDict(data []byte) []model.DictMatched {
	var dict_matches []model.DictMatched = nil

	matches := p.trie.Match(data)
	if len(matches) > 0 {
		for _, match := range matches {
			matched := model.DictMatched{
				Lable:    p.ColDicts[match.Pattern()].InfoType,
				Match:    string(match.Match()),
				Position: int(match.Pos()),
			}
			if dict_matches == nil {
				dict_matches = make([]model.DictMatched, 0)
			}
			dict_matches = append(dict_matches, matched)

			atomic.AddUint64(&p.ColDicts[match.Pattern()].CountMatch, 1)
		}

		p.CountMatchDict += uint64(len(matches))
	}

	return dict_matches
}
