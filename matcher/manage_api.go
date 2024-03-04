package matcher

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"data-matcher/utils"
)

type ManageApi struct {
	utils.RestapiHandler

	Host       string
	ApiPorts   []int
	Configfile string
	// Config        *MyConfig
	// Stats_stlog   *CompleteStatistic
	// Stats_apilog  *CompleteStatistic
	// Usercache     *UserCache
	// Accountcache  *AccountCache
	// Assetpolicy   *AssetPolicy
	// Authcache     *AuthCache
	// Sysdevcache   *SysdevCache
	// Holidaypolicy *HolidayPolicy
	// Chan_stlog_0  chan interface{}
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
	if len(p.ApiPorts) == 0 {
		return errors.New("not set api port")
	}
	localaddr := fmt.Sprintf(":%d", p.ApiPorts[0])

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
	log.Info("ManageApi listen https://", localaddr)
	if err := p.server.ListenAndServeTLS("", ""); err != nil {
		if err.Error() == "http: Server closed" {
			return nil
		}
		log.Error("ManageApi ListenAndServeTLS failed: ", err)
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
		return nil
	}
	tlscfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return tlscfg
}

var (
	CertPem = []byte(`-----BEGIN CERTIFICATE-----
	MIIDmzCCAoOgAwIBAgIJAKKqG9wnV0O3MA0GCSqGSIb3DQEBCwUAMGMxCzAJBgNV
	BAYTAkNOMQswCQYDVQQIDAJCSjELMAkGA1UEBwwCQkoxDTALBgNVBAoMBGlkc3Mx
	ETAPBgNVBAsMCHNlY3VyaXR5MRgwFgYDVQQDDA9pZG0uaWRzcy1jbi5jb20wIBcN
	MjQwMzAyMDYxODMzWhgPMjEyNDAyMDcwNjE4MzNaMGMxCzAJBgNVBAYTAkNOMQsw
	CQYDVQQIDAJCSjELMAkGA1UEBwwCQkoxDTALBgNVBAoMBGlkc3MxETAPBgNVBAsM
	CHNlY3VyaXR5MRgwFgYDVQQDDA9pZG0uaWRzcy1jbi5jb20wggEiMA0GCSqGSIb3
	DQEBAQUAA4IBDwAwggEKAoIBAQCyWSFjBCh5nWTagDQPMj14n9PwwbMp4SfRrB/X
	bJpfPsUbMXYaghUnl9sb7qf10pDW5zVh0c5M48vXoryebU4vGWxrcgeCwBI+TXRG
	VOq6/CJ3spM3zxtEg8vKEL81E4qYT2RItN5g+fWjeOAgxik9QS5eJrTBN4YaGd23
	y8XGV+rHMkbk3S13pS3XiJ4Qm8NYnPScwSY/Nt5BsULyBzCb58JzXaa4epEmopqz
	7/tU8RBcEyVW/UZ4tsJyPIfAGuCJUR4VIRELNhggaX2nF2lan40xKoOd3w55io6S
	Hi2E4ryhn248vBoqk77fkhR+mE1TfV1fJCZQPFI0zSibysPLAgMBAAGjUDBOMB0G
	A1UdDgQWBBRlbGT5LafVQc3tAekoLnoFH/bKjDAfBgNVHSMEGDAWgBRlbGT5LafV
	Qc3tAekoLnoFH/bKjDAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBO
	fvbdmBVHJmm1pPhdTVuSscr+43Ao8vHBgMASEwalZyeupj7MHS1USZR8v+gd1gUu
	FFsNCX3HA6SwQLxMsAR8HS3lCNexVOwfuI49vq0dped4kg7f8C+8c61NJKlteXXB
	HbQlrLq+sTmcKVtzgSwGEDjq9nlEr3HktPPRRA1h0DlWtrnriWgVPaUiDjZ+mrjV
	gvRnxNpai+UZ85sj8gVT8LROteWSLPWg8Soef8Ul7wUfSeAyeMIW9DSJgGnFkJ9B
	d0HmSCinjXn6SrBKC6kB3gP8pTiR611OyDOhwjXoHQL/WPNyWonKdbkjBIfX7xLf
	wRWnHaba5D096ufOAxEA
	-----END CERTIFICATE-----`)

	KeyPem = []byte(`-----BEGIN RSA PRIVATE KEY-----
	MIIEpQIBAAKCAQEAslkhYwQoeZ1k2oA0DzI9eJ/T8MGzKeEn0awf12yaXz7FGzF2
	GoIVJ5fbG+6n9dKQ1uc1YdHOTOPL16K8nm1OLxlsa3IHgsASPk10RlTquvwid7KT
	N88bRIPLyhC/NROKmE9kSLTeYPn1o3jgIMYpPUEuXia0wTeGGhndt8vFxlfqxzJG
	5N0td6Ut14ieEJvDWJz0nMEmPzbeQbFC8gcwm+fCc12muHqRJqKas+/7VPEQXBMl
	Vv1GeLbCcjyHwBrgiVEeFSERCzYYIGl9pxdpWp+NMSqDnd8OeYqOkh4thOK8oZ9u
	PLwaKpO+35IUfphNU31dXyQmUDxSNM0om8rDywIDAQABAoIBAQCd33JDcSHXDbGC
	DayHqyRpC6oT25MaRln2K5SAIH3CRBE80hrGulG5m530at05KGzYHxDNB2jD/X2q
	4z5uSznDTZEAx47IefdsOSntPCwQ2zIznNreszFjA/u4YfywIh00WErgZWLYm0uK
	qmxT9rX4qCNAaqjkxJ6rqivvD62BtUPWyPTdRq4OFyLTDqEdlQxaddYF8eQfBOpK
	h55lONeJyoc/ns43ftb33wtnWwc6XM9Lenttx5nFrtUTHDr++VMe0L5gQQ1V5vVt
	6QQ6Q5tf5B1ArfzlmL4bHBRC7f0iSNTLN9L8LwbJhSBsExhnFqYjc10gGYhB9MK7
	ld/gidPJAoGBAOhx+EQJi93NlFsaj2rXEy5jFpQH7IkDKWdhHmH/2EYNkZ6LSoqp
	hgyVUn5eaODBuVEq7SagUnZk4FeZo/U46KsHmucAJXIxJleApjz5n2324K7kaGaf
	XzMuQEEwb9RDnLi/29KBxlMr6FWJ82GQINVuHweNqK/nlr8wDBSZsApnAoGBAMRr
	ykjwmSHtNuo2pCojSbDdCHOMbu8fld8gtbfWSNMqVC0ntiHuN1/K54WHdy4JyQSS
	v1TdEl2ai2jbLRbYWuGkBr5bQmLUPD7Kh5Yl4bAZv+DVbRf4RAd1RdkeDJ+ezrSx
	o8YBQ8qaIO3ZnirUvEIv8nPvAoYN7sQ04cnJliT9AoGBAJWyAC7g7wBzCt35Ju+p
	fyLakYnX6I78SEfZldWLDN9gka1HC0RtlHS6HZxgdK56VDxfpsa/bRvuL0R7H8on
	UkAC79FgmL0HxieIJIcUQ4Zv/ZbkZg/hB1BQsvTImtxahq28cXcKOI0Ls96Srvjf
	9yU8fCNDKaXPQZfy+3Sw3Vx1AoGBAKMUz84BnVLSzk5l8aVeyRdEXXj6dzyon9mz
	Ic0x6CMTOPKIzyqay3UIVXPDRot96l2Wra77If1/jBISL/yQw9wmQMcZpCPEDQUh
	SLO8XgbFSk+VRE+rfGgo0UZ0MYzx4LOb7ds/P5beo0p37V+oY2ocvxPMtO6ycLSN
	J45Phg7NAoGALxuGpbxYX7cA060qyyNUHCOYMt7yPT7L5ex1KCmcOdyr2OkazSw2
	Kj6ziMf/7Qk+uIHjVlecupxJ/NjFRrWUFRQdNdEQ7uQ/hSCu7ITBW/1F6vE7sLbJ
	DVtbIZqc+VvNUpq0HrAnWayU0voxfcSwuDg3yvWK8fR8Bg57hGRnCPw=
	-----END RSA PRIVATE KEY-----`)
)
