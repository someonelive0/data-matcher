package matcher

import (
	"fmt"

	"data-matcher/model"
)

func (p *Worker) processFlowDns(flowDns *model.FlowDns) error {
	if len(flowDns.Dns.Rrname) == 0 {
		return fmt.Errorf("rrname is empty")
	}

	return nil
}
