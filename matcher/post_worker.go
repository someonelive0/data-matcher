package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"data-matcher/model"
)

type PostWorker struct {
	Httpch chan *model.MsgHttp `json:"-"`
	AppMap sync.Map            `json:"-"` // app, number
	ApiMap sync.Map            `json:"-"` // api, number
	IpMap  sync.Map            `json:"-"` // client ip, number
}

func (p *PostWorker) Run() {
	for msgHttp := range p.Httpch {
		p.discoverApp(msgHttp)

		p.discoverApi(msgHttp)

		p.discoverIp(msgHttp)

		p.discoverAccount(msgHttp)

		p.disposeToken(msgHttp)

		p.disposeUsername(msgHttp)
	}
}

func (p *PostWorker) Stop() {
}

// discover app from msg of subject flow.http
func (p *PostWorker) discoverApp(msgHttp *model.MsgHttp) error {
	app := fmt.Sprintf("%s/%s/%d", msgHttp.Http.Hostname, msgHttp.Dest_ip, msgHttp.Dest_port)
	_, ok := p.AppMap.Load(app)
	if !ok {
		p.AppMap.Store(app, 1)
	}

	return nil
}

// discover api from msg of subject flow.http
func (p *PostWorker) discoverApi(msgHttp *model.MsgHttp) error {
	api := msgHttp.Http.Url
	i := strings.IndexByte(api, '?') // erase ?
	if i != -1 {
		api = api[:i]
	}
	i = strings.IndexByte(api, '#') // erase #
	if i != -1 {
		api = api[:i]
	}

	api = msgHttp.Http.Hostname + api
	_, ok := p.ApiMap.Load(api)
	if !ok {
		p.ApiMap.Store(api, 1)
	}
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
	if len(msgHttp.Src_ip) == 0 {
		return nil
	}

	_, ok := p.IpMap.Load(msgHttp.Src_ip)
	if !ok {
		p.IpMap.Store(msgHttp.Src_ip, 1)
	}

	return nil
}

func (p *PostWorker) Dump() []byte {
	b, _ := json.Marshal(p)
	return b
}
