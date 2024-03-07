package matcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"

	"data-matcher/utils"
)

type NatsConfig struct {
	Servers  []string `yaml:"servers"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password" json:"-"`
}

type Flow struct {
	QueueName string `yaml:"queue_name"`
	Subject   string `yaml:"subject"`
}

type MyConfig struct {
	Filename string    `yaml:"filename" json:"filename" xml:"filename,attr"`
	LoadTime time.Time `yaml:"load_time" json:"load_time" xml:"load_time,attr"`

	Version     string `yaml:"version" json:"version"`
	Host        string `yaml:"host" json:"host"`
	ManagePort  int    `yaml:"manage_port" json:"manage_port" `
	ChannelSize int    `yaml:"channel_size" json:"channel_size"`
	Workers     int    `yaml:"workers" json:"workers"`
	RulesFile   string `yaml:"rules_file" json:"rules_file"`

	NatsConfig NatsConfig `yaml:"nats"`
	HttpFlow   Flow       `yaml:"http_flow"`

	Statsviz bool `yaml:"statsviz"`
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

	cpus := runtime.NumCPU()
	if myconfig.Workers <= 0 || myconfig.Workers >= cpus {
		fmt.Printf("workers [%d] is less-eque 0 or bigger than cpu number, set it to cpu number - 1\n", myconfig.Workers)
		if cpus == 1 {
			myconfig.Workers = 1
		} else {
			myconfig.Workers = cpus - 1
		}
	}

	// 补上rules_file的路径
	myconfig.RulesFile = path.Dir(myconfig.Filename) + "/" + myconfig.RulesFile

	return myconfig, nil
}
