package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"data-matcher/utils"
)

var (
	arg_server   = flag.String("server", "nats://localhost:4222", "nats server hostname or ip")
	arg_user     = flag.String("user", "", "nats user")
	arg_password = flag.String("password", "", "")
	arg_subject  = flag.String("subject", "flow", "subject of nats to receive message")
	arg_num      = flag.Int("num", 10000, "message numbers, default 10000")
	arg_file     = flag.String("file", "", "filename of messages with json in one line")
	msg          = `{"timestamp":"2024-01-09T16:36:00.089895+0800","response_time":0,"flow_id":696134833031826,"in_iface":"ens37,ens33","event_type":"http","pcap_filename":"","src_mac":"e8:65:49:4b:46:55","src_ip":"10.10.20.113","src_port":56250,"dst_mac":"00:27:0d:ff:5b:e1","dest_ip":"36.25.240.158","dest_port":9911,"proto":"TCP","tx_id":0,"http":{"hostname":"www.jfdaily.com","url":"\/Wsusadmin\/Errors\/BrowserSettings.aspx","http_user_agent":"Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)","http_content_type":"text\/html","http_method":"GET","protocol":"HTTP\/1.1","status":400,"length":666,"request_headers":[{"name":"Connection","value":"Keep-Alive"},{"name":"Host","value":"www.jfdaily.com"},{"name":"Pragma","value":"no-cache"},{"name":"User-Agent","value":"Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)"},{"name":"Accept","value":"image\/gif, image\/x-xbitmap, image\/jpeg, image\/pjpeg, image\/png, *\/*"},{"name":"Accept-Language","value":"en"},{"name":"Accept-Charset","value":"iso-8859-1,*,utf-8"}],"response_headers":[{"name":"Server","value":"nginx"},{"name":"Date","value":"Fri, 18 Oct 2019 09:13:22 GMT"},{"name":"Content-Type","value":"text\/html"},{"name":"Content-Length","value":"666"},{"name":"Connection","value":"close"}],"reqHeaders":"Connection: Keep-Alive\r\nHost: www.jfdaily.com\r\nPragma: no-cache\r\nUser-Agent: Mozilla\/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident\/4.0)\r\nAccept: image\/gif, image\/x-xbitmap, image\/jpeg, image\/pjpeg, image\/png, *\/*\r\nAccept-Language: en\r\nAccept-Charset: iso-8859-1,*,utf-8\r\n\r\n","respHeaders":"Server: nginx\r\nDate: Fri, 18 Oct 2019 09:13:22 GMT\r\nContent-Type: text\/html\r\nContent-Length: 666\r\nConnection: close\r\n\r\n","respBody":"<html>\r\n<head><title>400 The plain HTTP request was sent to HTTPS port<\/title><\/head>\r\n<body bgcolor=\"white\">\r\n<center><h1>400 Bad Request<\/h1><\/center>\r\n<center>The plain HTTP request was sent to HTTPS port<\/center>\r\n<hr><center>nginx 36.25.240.158<\/center>\r\n<\/body>\r\n<\/html>\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n","respLen":666,"reqBody":"","reqLen":0,"totalLen":1063},"host":"172.16.0.230"}`
)

func init() {
	flag.Parse()
}

func main() {
	nc, err := utils.NatsConnect(*arg_server, *arg_user, *arg_password)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nc.Close()
	log.Printf("connect %s success by user %s\n", *arg_server, *arg_user)

	t0 := time.Now()
	n := 0
	// 如果数据参数为空，则循环发送单样本数据
	if len(*arg_file) > 0 {
		n, err = sendFileMsg(nc)
	} else {
		n, err = sendSingleMsg(nc)
	}
	if err != nil {
		return
	}

	if err = nc.Flush(); err != nil {
		log.Println("flush msg failed: ", err)
	}
	log.Printf("end sent %d messages in %v\n", n, time.Now().Sub(t0))

	// time.Sleep(5 * time.Second)
}

func sendSingleMsg(nc *nats.Conn) (int, error) {
	log.Printf("sending single message loop %d with length %d\n", *arg_num, len(msg))
	for i := 0; i < *arg_num; i++ {
		if err := utils.NatsSendMsg(nc, *arg_subject, msg); err != nil {
			log.Println("sendmsg failed: ", err)
			return 0, err
		}

		if i%100 == 0 {
			// if err = nc.FlushTimeout(1 * time.Second); err != nil {
			// 	log.Println("FlushTimeout msg failed: ", err)
			// }
			if err := nc.Flush(); err != nil {
				log.Println("flush msg failed: ", err)
				return i, err
			}
		}
	}

	return *arg_num, nil
}

func sendFileMsg(nc *nats.Conn) (int, error) {
	fp, err := os.Open(*arg_file)
	if err != nil {
		log.Printf("open data file [%s] failed: %s\n", *arg_file, err)
		return 0, err
	}
	defer fp.Close()

	i := 0
	m := make(map[string]interface{})
	submap := make(map[string]int)

	scanner := bufio.NewScanner(fp)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &m)
		if err != nil {
			fmt.Printf("Umarshal failed: %s\n%s\n", err, line)
			continue
		}
		if event_type, ok := m["event_type"]; ok {
			// log.Printf("event_type: %s\n", event_type)

			subject := *arg_subject + "." + event_type.(string) // subject is like flow.http
			if err := nc.Publish(subject, line); err != nil {
				log.Println("sendmsg failed: ", err)
				return 0, err
			}
			i++

			submap[subject]++
		}
	}

	log.Printf("subjects: %#v", submap)
	if scanner.Err() != nil {
		log.Println("file scanner failed: ", scanner.Err())
	}

	return i, scanner.Err()
}
