package matcher

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

type DnsItem struct {
	Rrname string
	Value  string
}

func (p *Worker) processDns(m *nats.Msg) (*DnsItem, error) {
	dns, err := sonic.Get(m.Data, "dns")
	if err != nil {
		return nil, err
	}
	rrname, err := dns.GetByPath("rrname").String()
	if err != nil {
		return nil, err
	}
	if len(rrname) == 0 {
		return nil, fmt.Errorf("rrname is empty")
	}
	value, err := dns.Raw()
	if err != nil {
		return nil, err
	}

	dnsitem := &DnsItem{
		Rrname: rrname,
		Value:  value,
	}

	return dnsitem, nil
}
