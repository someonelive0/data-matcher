package model

/*
 来自 Kafka 的消息样本，这些消息有流转探针turboflow产生，写入kafka

 具体样本参考 flow_msg_sample.md
*/

/*
 每一个消息都具有消息头，但是缺个版本描述，
 event_type 用来区分消息头之后的具体内容，event_type包括： http, dns, ssh, rdp, elasticsearch, pop3, stmp
*/
type FlowHeader struct {
	Timestamp     string `json:"timestamp"`
	Response_time int    `json:"response_time"`
	Flow_id       int    `json:"flow_id"`
	Event_type    string `json:"event_type"` // event_type表示消息头之后的内容，包括： http, dns, ssh, rdp, elasticsearch, pop3, stmp
	Src_mac       string `json:"src_mac"`
	Src_ip        string `json:"src_ip"`
	Src_port      int    `json:"src_port"`
	Dst_mac       string `json:"dst_mac"`
	Dest_ip       string `json:"dest_ip"`
	Dest_port     int    `json:"dest_port"`
	Proto         string `json:"proto"`
	Tx_id         int    `json:"tx_id"`
	Host          string `json:"host"`
}

/*
 当event_type=http时，后续内容为http消息
*/
type FlowHttp struct {
	FlowHeader
	Http HttpMsg `json:"http"`
}

/*
 http消息在 消息头之后，用"http" JSON标签表示
*/
type HttpMsg struct {
	Hostname          string       `json:"hostname"`
	Url               string       `json:"url"`
	Http_user_agent   string       `json:"http_user_agent"`
	Http_content_type string       `json:"http_content_type"`
	Http_method       string       `json:"http_method"`
	Protocol          string       `json:"protocol"`
	Status            int          `json:"status"`
	Length            int          `json:"length"`
	Request_headers   []HttpHeader `json:"request_headers"`
	Response_headers  []HttpHeader `json:"response_headers"`
	ReqHeaders        string       `json:"reqHeaders"`
	ReqBody           string       `json:"reqBody"`
	ReqLen            int          `json:"reqLen"`
	RespHeaders       string       `json:"respHeaders"`
	RespBody          string       `json:"respBody"`
	RespLen           int          `json:"respLen"`
	TotalLen          int          `json:"totalLen"`

	// 匹配正则和字典的结果
	Value_regex []RegexMatched `json:"value_regex"`
	Column_dict []DictMatched  `json:"column_dict"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RegexMatched struct {
	Lable    string `json:"lable"`
	Position []int  `json:"position"`
}

type DictMatched struct {
	Lable    string `json:"lable"`
	Match    string `json:"match"`
	Position int    `json:"position"`
}

/*
 当event_type=http时，后续内容为http消息
*/
type FlowDns struct {
	FlowHeader
	Dns DnsMsg `json:"dns"`
}

/*
 dns消息在 消息头之后，用"dns" JSON标签表示
*/
type DnsMsg struct {
	Version int                 `json:"version"`
	Type    string              `json:"type"`
	Dd      int                 `json:"id"`
	Flags   string              `json:"flags"`
	Qr      bool                `json:"qr"`
	Rd      bool                `json:"rd"`
	Ra      bool                `json:"ra"`
	Rrname  string              `json:"rrname"`
	Rrtype  string              `json:"rrtype"`
	Rcode   string              `json:"rcode"`
	Answers []DnsAnswer         `json:"answers"`
	Grouped map[string][]string `json:"grouped"`
}

type DnsAnswer struct {
	Rrname string `json:"rrname"`
	Rrtype string `json:"rrtype"`
	Ttl    int    `json:"ttl"`
	Rdata  string `json:"rdata"`
}
