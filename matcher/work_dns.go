package matcher

import (
	"fmt"

	"data-matcher/model"
)

func (p *Worker) processMsgDns(msgDns *model.MsgDns) error {
	if len(msgDns.Dns.Rrname) == 0 {
		return fmt.Errorf("rrname is empty")
	}

	return nil
}
