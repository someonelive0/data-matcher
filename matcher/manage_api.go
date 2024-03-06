package matcher

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

type ManageApi struct {
	utils.RestapiHandler

	// Host       string
	Port   int
	Config *MyConfig
	// Stats_stlog   *CompleteStatistic
	// Stats_apilog  *CompleteStatistic
	// Usercache     *UserCache
	// Accountcache  *AccountCache
	// Assetpolicy   *AssetPolicy
	// Authcache     *AuthCache
	// Sysdevcache   *SysdevCache
	// Holidaypolicy *HolidayPolicy
	MsgChan chan *nats.Msg
	// Chan_stlog_1  chan interface{}
	// Chan_apilog_0 chan interface{}
	// Chan_apilog_1 chan interface{}
	// Input         *MsgInput
	// Stlogoutput   *MylogOutput
	// Apilogoutput  *MylogOutput
	Myerrors *utils.MyErrors

	server *http.Server
}

func (p *ManageApi) Run() error {
	localaddr := fmt.Sprintf(":%d", p.Port)

	r := mux.NewRouter().StrictSlash(true)
	r.Use(utils.AuthMiddleware)
	r.HandleFunc("/", p.StatusHandler).Methods("GET")
	r.HandleFunc("/status", p.StatusHandler).Methods("GET")
	r.HandleFunc("/debug", p.DebugHandler).Methods("GET", "PUT")
	r.HandleFunc("/version", p.VersionHandler).Methods("GET")
	r.HandleFunc("/changelog", p.ChangelogHandler).Methods("GET")
	r.HandleFunc("/log", p.LogHandler).Methods("GET")

	r.HandleFunc("/dump", p.dumpHandler).Methods("GET")
	r.HandleFunc("/statistic", p.statisticHandler).Methods("GET")
	r.HandleFunc("/errors", p.errorsHandler).Methods("GET")
	r.HandleFunc("/config", p.configHandler).Methods("GET", "POST")
	r.HandleFunc("/config/", p.configHandler).Methods("GET", "POST")
	r.HandleFunc("/config/{title:[a-z0-9_\\-]+}", p.configItemHandler).Methods("GET", "POST")

	r.HandleFunc("/policy", p.UploadHandler).Methods("POST")
	r.HandleFunc("/policy", p.policyReloadHandler).Methods("PUT")
	r.HandleFunc("/policy/", p.policyHandler).Methods("GET")
	r.HandleFunc("/policy/{title:[a-z0-9_\\-]+}", p.policyHandler).Methods("GET")

	p.server = &http.Server{
		Addr:         localaddr,
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		TLSConfig:    LoadTLSCert(),
		Handler:      r,
	}

	//defer log.Debug("ManageApi closed")
	log.Infof("ManageApi listen https://%s", localaddr)
	if err := p.server.ListenAndServeTLS("", ""); err != nil {
		if err.Error() == "http: Server closed" {
			return nil
		}
		log.Errorf("ManageApi ListenAndServeTLS failed: %s", err)
		return err
	}
	return nil
}

func (p *ManageApi) Stop() error {
	if p.server != nil {
		p.server.Shutdown(context.Background())
		p.server = nil
	}
	return nil
}

/*
TLS and HTTPS cert, self signed cert and key pem file with PKCS#1-PEM

openssl genrsa -out key.pem 2048
openssl req -new -x509 -sha256 -key key.pem -out cert.pem -days 36500 -subj "/C=CN/ST=BJ/L=BJ/O=idss/OU=security/CN=idm.idss-cn.com"
openssl rsa -in key.pem -noout -text
openssl x509 -in cert.pem -noout -text
*/
func LoadTLSCert() *tls.Config {
	//cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	cert, err := tls.X509KeyPair(CertPem, KeyPem)
	if err != nil {
		log.Errorf("tls.X509KeyPair failed: %s", err)
		return nil
	}
	tlscfg := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
	return tlscfg
}

var (
	CertPem = []byte(`-----BEGIN CERTIFICATE-----
MIIDmzCCAoOgAwIBAgIJAK5qKMsSAwwzMA0GCSqGSIb3DQEBCwUAMGMxCzAJBgNV
BAYTAkNOMQswCQYDVQQIDAJCSjELMAkGA1UEBwwCQkoxDTALBgNVBAoMBGlkc3Mx
ETAPBgNVBAsMCHNlY3VyaXR5MRgwFgYDVQQDDA9pZG0uaWRzcy1jbi5jb20wIBcN
MjQwMzA0MDkwMDIxWhgPMjEyNDAyMDkwOTAwMjFaMGMxCzAJBgNVBAYTAkNOMQsw
CQYDVQQIDAJCSjELMAkGA1UEBwwCQkoxDTALBgNVBAoMBGlkc3MxETAPBgNVBAsM
CHNlY3VyaXR5MRgwFgYDVQQDDA9pZG0uaWRzcy1jbi5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDhd7/M0jl7OsLns5k21XnUdbWVlehTi/GR1XDn
F9UsRS9n2Hz8bXk8XagHX65wG3u1tjulOdZNe+2tSJ8pyUC2D0Z9Rbg/EDL+fqlr
Got8gfqatfZ1sD/GbtINqq8z8L4v7WWclExIQDw81Vc5gFVD3G8iUnC8uIG8wTKx
gB+EJ0nSpK9NpMy4O4CIOrmlCaivVR6Uqjb0HOLhDmBXR858zH0o1h5I6Iir696z
WXu6PYaZG3Tar8VS6WGA+7ybzpZJVBJOKBKrzMW0Xj3KbYKfQJqD86USKEDzubHx
GeWWvRUIL98FPFFJslbQ0N2+EvbxRPydxffmrVV6pieSWpw1AgMBAAGjUDBOMB0G
A1UdDgQWBBRFoK4eWcr2pij5kdAI674/dyG2SzAfBgNVHSMEGDAWgBRFoK4eWcr2
pij5kdAI674/dyG2SzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCe
RaTCKVFhzp7zXXh7tG7LXgHhmseSs+boz6nOTPOdG2e9CkjkeXfy7hDYpTdE5pys
qweKQ3NhyubxkGVJLlSeLqqiA7XbcuwCh/zgdnaUMzDBJvy5PYHkLxn1nGBI68fX
03N1nA9y+4k0luYT3P40P5cgBkR7MhOLYL16v5Uh4mt4CdoL/89YX+yIIfCnqTUf
gN2E2TEbyYAfhvmu/q9BREhiPOBK0RhldAsmyTPTizEC6R6g0oHM4nlFzayVB2qA
7j6CHBXa+H19B16Hw6JxXGz6KEXa2pKHehHu3UzDrbiJIWly5yH1fa74MDLAnrGq
KbMxgtkmhNMoK5yMp3Li
-----END CERTIFICATE-----
`)

	KeyPem = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA4Xe/zNI5ezrC57OZNtV51HW1lZXoU4vxkdVw5xfVLEUvZ9h8
/G15PF2oB1+ucBt7tbY7pTnWTXvtrUifKclAtg9GfUW4PxAy/n6paxqLfIH6mrX2
dbA/xm7SDaqvM/C+L+1lnJRMSEA8PNVXOYBVQ9xvIlJwvLiBvMEysYAfhCdJ0qSv
TaTMuDuAiDq5pQmor1UelKo29Bzi4Q5gV0fOfMx9KNYeSOiIq+ves1l7uj2GmRt0
2q/FUulhgPu8m86WSVQSTigSq8zFtF49ym2Cn0Cag/OlEihA87mx8Rnllr0VCC/f
BTxRSbJW0NDdvhL28UT8ncX35q1VeqYnklqcNQIDAQABAoIBAQDgx5rvHxLxeQbB
KrtwAGniV6u9wuMJD/a5Flrl+UusRBlb5WfN3XJFrXWMTGbDG5M8+L4EHmI5g3jU
Dhen+B1MpHP5Bl4GeSbts+dBgQhZC9iFDy3z7M/YC7ncqDLdjIB/laR88xgN3ARx
/ZtlFz1qV0RwqlH8w9GMFextK7e/7fAymM/S6m8+drdhj6EWPESQwvORNizi8OK2
y/TpYW0yPXEGzRVi0kJdZDi9ipcw1Lsec+UtJag4UEiHHyOA7VASwLEJ2Kw+HJGu
YxXLAQXsTehm2yuB+cdV+4xgwnLTDvKupCrf2fjViN90zz8q5qyBhpI3eeHUJZb+
ZQJkC0XhAoGBAPWVWblY6Fu5cM8+cpzF2rWArhKccTzN040WeYvP5zDmlkzcULyu
jxsWgGXVem+d4nBz8XirBFVWKOjKuuxyAMdLv3aCivGCh8SsSHsGwwExK7cPrWB4
TCqwKqdJnUnpMj/T2C1CYFrY3Qt5EX0T0b+nWCkVZdHYluRNJN0r6DmtAoGBAOsH
+YE3QuVaV+Y1pHJ/4YMvkDzpGntYmayWXpKk0ysOhBGSZJXsHVewwCGLOnly7XX6
ecmXTCiOQRLWqFYCyKj8+V9VE0XpDUXqrSJAzfcRMtsK+NcMjfU4HEIRQP/B1NT8
lYBKYndPgK2TbRCz8ukC4IU5dLO6L83Ae4ok8c2pAoGAG1FDq8RiBGH6VHND2ICB
tZLcyiEwz2ytzZHkb1LvCpd7vIz9Rh+8t2ynV6yJdAUB/TRIdf2/+6Yb4tk6Nbbw
szqPz6Txw6+bXpszbMvxwR4xGKnbxVFcV5tFA1rC7kfMWSE9eLtbcH+TBwWullUw
DbuVqOxCaTdIgZi7MwcBS/kCgYBAewOcy8hynAKZigX/083O6/GqhFlblcczbl2r
5cR5f5YELCGkcA7sy/UqPsRgJYO4ZmubPwMJ7V01CedNEZ5znlPcL78F4xZdJDEz
wIvBSNqm9a+ncC5SJH68MXefs1HszQ9HDyFMkmc/N78oYfY2ry9h3Y2C8YXD8Rbz
o4cjYQKBgQDAj0TglE/pNlf1TC2h7Eigm87Shat2ILV3vJLPDHnlY03/FA6uNSaf
XWDTMrfRKdsWNe3NFo8o1Z72yyYTrPmgCtSAbuo2+bT9Qu1ccHck4zYVzpskGwJe
+HvPAG7YzAMZSCKOilCQfSojmnb1euMhfc4Hgbwu3lOG1ilpjv0hWA==
-----END RSA PRIVATE KEY-----
`)
)