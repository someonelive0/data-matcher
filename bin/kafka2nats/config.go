package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"data-matcher/utils"
)

// SASL machanism, such as PLAIN, SCRAM
type KafkaConfig struct {
	Brokers   []string `yaml:"brokers" json:"brokers"`
	Mechanism string   `yaml:"mechanism" json:"mechanism"`
	User      string   `yaml:"user" json:"user"`
	Password  string   `yaml:"password" json:"-"`
	Group     string   `yaml:"group" json:"group"`
}

type NatsConfig struct {
	Servers  []string `yaml:"servers"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password" json:"-"`
}

type BridgeConfig struct {
	Topic     string `yaml:"topic"`
	Subject   string `yaml:"subject"`
	Jetstream bool   `yaml:"jetstream"`
}

type MyConfig struct {
	Filename string    `yaml:"filename" json:"filename" xml:"filename,attr"`
	LoadTime time.Time `yaml:"load_time" json:"load_time" xml:"load_time,attr"`

	Version    string `yaml:"version" json:"version"`
	Host       string `yaml:"host" json:"host"`
	ManagePort int    `yaml:"manage_port" json:"manage_port" `

	KafkaConfig KafkaConfig `yaml:"kafka"`
	NatsConfig  NatsConfig  `yaml:"nats"`

	Kafka2Nats BridgeConfig `yaml:"kakfa2nats"`
	Nats2Kafka BridgeConfig `yaml:"nats2kafka"`
}

func (p *MyConfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", " ")
	return b
}

func LoadConfig(filename string) (*MyConfig, error) {
	// 都配置文件，如果文件不存在则从模块文件tpl复制成配置文件。思路是考虑到不覆盖已有现场配置文件。
	if !utils.ExistedOrCopy(filename, filename+".tpl") {
		return nil, fmt.Errorf("config file [%s] or template file are not found", filename)
	}

	// 打开 YAML 文件
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	// 创建解析器
	decoder := yaml.NewDecoder(fp)

	// 配置对象
	myconfig := &MyConfig{
		Filename: filename,
		LoadTime: time.Now(),
	}

	// 解析 YAML 数据
	err = decoder.Decode(myconfig)
	if err != nil {
		// fmt.Println("Error decoding YAML:", err)
		return nil, err
	}

	return myconfig, nil
}
