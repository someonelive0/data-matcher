package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

// process msg with subject flow.dns
func (p *Worker) proccessDns(m *nats.Msg) (map[string]interface{}, error) {
	jsonmap := make(map[string]interface{})
	err := json.Unmarshal(m.Data, &jsonmap) // TODO，改成sonic解析
	if err != nil {
		log.Errorf("worker unmarshal dns failed: %s\n", err)
		return nil, err
	}

	if dnsmap, ok := jsonmap["dns"]; ok {
		// b, _ := json.Marshal(dnsmap)
		// log.Debugf("DNS: %s", b)
		return dnsmap.(map[string]interface{}), nil
	}

	return nil, fmt.Errorf("not found dns in json")
}
