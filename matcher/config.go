package matcher

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"data-matcher/utils"
)

type MyConfig struct {
	Filename string    `yaml:"filename" json:"filename" xml:"filename,attr"`
	LoadTime time.Time `yaml:"load_time" json:"load_time" xml:"load_time,attr"`

	Version string `yaml:"version" json:"version"`
	Host    string `yaml:"host" json:"host"`

	Service struct {
		Port int `yaml:"port" json:"port" `
	} `yaml:"service"`

	// Kafka KafkaConfig `yaml:"kafka" json:"kafka"`
	//Nacos      utils.NacosConfig `yaml:"nacos" json:"nacos"`
}

func (p *MyConfig) Dump() []byte {
	b, _ := json.MarshalIndent(p, "", " ")
	return b
}

func LoadConfig(filename string) (*MyConfig, error) {
	if !utils.ExistedOrCopy(filename, filename+".tpl") {
		return nil, fmt.Errorf("config file [%s] or template file are not found", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config file [%s] open failed: %s", filename, err)
	}

	myconfig := &MyConfig{
		Filename: filename,
		LoadTime: time.Now(),
	}
	err = yaml.Unmarshal(data, myconfig)
	if err != nil {
		return nil, fmt.Errorf("config file [%s] unmarshal yaml failed: %s", filename, err)
	}

	return myconfig, nil
}
