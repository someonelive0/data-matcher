package matcher

import (
	"data-matcher/model"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type PostWorker struct {
	Httpch chan *model.MsgHttp `json:"-"`
}

func (p *PostWorker) Run() {
	for msgHttp := range p.Httpch {
		p.discoverApp(msgHttp)
	}
}

func (p *PostWorker) Stop() {
}

// discover app from msg of subject flow.http
func (p *PostWorker) discoverApp(msgHttp *model.MsgHttp) error {

	s := fmt.Sprintf("%s/%s/%d", msgHttp.Http.Hostname, msgHttp.Dest_ip, msgHttp.Dest_port)
	log.Debugf("---> %s", s)

	return nil
}

// discover api from msg of subject flow.http
func (p *PostWorker) discoverApi(msgHttp *model.MsgHttp) error {
	return nil
}

// dispose username from msg of subject flow.http
func (p *PostWorker) disposeUsername(msgHttp *model.MsgHttp) error {
	return nil
}

// dispose token from msg of subject flow.http
func (p *PostWorker) disposeToken(msgHttp *model.MsgHttp) error {
	return nil
}

// discover account from msg of subject flow.http
func (p *PostWorker) discoverAccount(msgHttp *model.MsgHttp) error {
	return nil
}

// discover client ip from msg of subject flow.http
func (p *PostWorker) discoverIp(msgHttp *model.MsgHttp) error {
	return nil
}
