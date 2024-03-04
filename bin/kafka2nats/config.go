package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
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
	KafkaConfig KafkaConfig `yaml:"kafka"`
	NatsConfig  NatsConfig  `yaml:"nats"`

	Kafka2Nats BridgeConfig `yaml:"kakfa2nats"`
	Nats2Kafka BridgeConfig `yaml:"nats2kafka"`
}

func ReadConfig(filename string) (*MyConfig, error) {
	// 打开 YAML 文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建解析器
	decoder := yaml.NewDecoder(file)

	// 配置对象
	var config MyConfig

	// 解析 YAML 数据
	err = decoder.Decode(&config)
	if err != nil {
		// fmt.Println("Error decoding YAML:", err)
		return nil, err
	}

	return &config, nil
}

func DumpConfig(myconfig *MyConfig) []byte {
	b, err := json.Marshal(*myconfig)
	if err != nil {
		return nil
	}
	return b
}
