package engine

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"data-matcher/utils"
)

type MaskRuleItem struct {
	RuleName      string   `yaml:"RuleName"`
	MaskType      string   `yaml:"MaskType"` // one of [CHAR, TAG, REPLACE, EMPTY, ALGO ]
	Value         string   `yaml:"Value"`
	Offset        int32    `yaml:"Offset"`
	Padding       int32    `yaml:"Padding"`
	Length        int32    `yaml:"Length"`
	Reverse       bool     `yaml:"Reverse"`
	IgnoreCharSet string   `yaml:"IgnoreCharSet"`
	IgnoreKind    []string `yaml:"IgnoreKind"` // one of [NUMERIC, ALPHA_UPPER_CASE, ALPHA_LOWER_CASE, WHITESPACE, PUNCTUATION]
}

type RuleItem struct {
	RuleID      int32  `yaml:"RuleID"`
	InfoType    string `yaml:"InfoType"`
	Description string `yaml:"Description"`
	EnName      string `yaml:"EnName"`
	CnName      string `yaml:"CnName"`
	Level       string `yaml:"Level"` // L1 (least Sensitive) ~ L4 (Most Sensitive)
	// (KReg || KDict) && (VReg || VDict)
	Detect struct {
		KReg  []string `yaml:"KReg"`       // Regex List for Key
		KDict []string `yaml:"KDict,flow"` // Dict for Key
		VReg  []string `yaml:"VReg"`       // Regex List for Value
		VDict []string `yaml:"VDict,flow"` // Dict for Value
	} `yaml:"Detect"`
	// result which is hit by blacklist will not returned to caller
	Filter struct {
		// BReg || BDict
		BReg  []string `yaml:"BReg"`       // Regex List for BlackList
		BDict []string `yaml:"BDict,flow"` // Dict for BlackList
		BAlgo []string `yaml:"BAlgo"`      // Algorithm List for BlackList, one of [ MASKED ]
	} `yaml:"Filter"`
	// result need pass verify process before retured to caller
	Verify struct {
		// CReg || CDict
		CReg  []string `yaml:"CReg"`       // Regex List for Context Verification
		CDict []string `yaml:"CDict,flow"` // Dict for Context Verification
		VAlgo []string `yaml:"VAlgo"`      // Algorithm List for Verification, one of [ IDVerif , CardVefif ]
	} `yaml:"Verify"`
	Mask    string            `yaml:"Mask"` // MaskRuleItem.RuleName for Mask
	ExtInfo map[string]string `yaml:"ExtInfo"`
}

type RulesConfig struct {
	Global struct {
		Date    string `yaml:"Date"`
		Version string `yaml:"Version"`
	} `yaml:"Global"`
	MaskRules []MaskRuleItem `yaml:"MaskRules"`
	Rules     []RuleItem     `yaml:"Rules"`
}

// newDlpConfImpl implements newDlpConf by receving conf content string
func NewRulesConfig(filename string) (*RulesConfig, error) {
	// 都配置文件，如果文件不存在则从模块文件tpl复制成配置文件。思路是考虑到不覆盖已有现场配置文件。
	if !utils.ExistedOrCopy(filename, filename+".tpl") {
		return nil, fmt.Errorf("config file [%s] or template file are not found", filename)
	}

	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &RulesConfig{}
	if err := yaml.Unmarshal(bs, config); err == nil {
		if err := config.Verify(); err == nil {
			return config, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (p *RulesConfig) GetValueReg() []string {
	ss := make([]string, 0)
	for _, rule := range p.Rules {
		if len(rule.Detect.VReg) != 0 {
			ss = append(ss, rule.Detect.VReg...)
		}
	}
	return ss
}

func (p *RulesConfig) GetColDict() []string {
	ss := make([]string, 0)
	for _, rule := range p.Rules {
		if len(rule.Verify.CDict) != 0 {
			ss = append(ss, rule.Verify.CDict...)
		}
	}
	return ss
}

var (
	defMaskTypeSet []string = []string{"CHAR", "TAG", "REPLACE", "ALGO"}
	defMaskAlgo    []string = []string{"BASE64", "MD5", "CRC32", "ADDRESS", "NUMBER", "DEIDENTIFY"}
	defIgnoreKind  []string = []string{"NUMERIC", "ALPHA_UPPER_CASE", "ALPHA_LOWER_CASE", "WHITESPACE", "PUNCTUATION"}
)

func (I *RulesConfig) Verify() error {
	// Global

	// MaskRules
	for _, rule := range I.MaskRules {
		// MaskType
		if inList(rule.MaskType, defMaskTypeSet) == -1 {
			return fmt.Errorf("%s, Mask RuleName:%s, MaskType:%s is not suppored", "ERR_CONF_VERIFY_FAILED", rule.RuleName, rule.MaskType)
		}
		if strings.Compare(rule.MaskType, "ALGO") == 0 {
			if inList(rule.Value, defMaskAlgo) == -1 {
				return fmt.Errorf("%s, Mask RuleName:%s, ALGO Value: %s is not supported", "ERR_CONF_VERIFY_FAILED", rule.RuleName, rule.Value)
			}
		}
		if !(rule.Offset >= 0) {
			return fmt.Errorf("%s, Mask RuleName:%s, Offset: %d need >=0", "ERR_CONF_VERIFY_FAILED", rule.RuleName, rule.Offset)
		}
		if !(rule.Length >= 0) {
			return fmt.Errorf("%s, Mask RuleName:%s, Length: %d need >=0", "ERR_CONF_VERIFY_FAILED", rule.RuleName, rule.Length)
		}
		for _, kind := range rule.IgnoreKind {
			if inList(kind, defIgnoreKind) == -1 {
				return fmt.Errorf("%s, Mask RuleName:%s, IgnoreKind: %s is not supported", "ERR_CONF_VERIFY_FAILED", rule.RuleName, kind)
			}
		}
	}
	// Rules
	for _, rule := range I.Rules {
		de := rule.Detect
		// at least one detect rule
		if len(de.KReg) == 0 && len(de.KDict) == 0 && len(de.VReg) == 0 && len(de.VDict) == 0 {
			return fmt.Errorf("%s, RuleID:%d, Detect field missing", "ERR_CONF_VERIFY_FAILED", rule.RuleID)
		}
	}
	return nil
}

// inList finds item in list
func inList(item string, list []string) int {
	for i, v := range list {
		if strings.Compare(item, v) == 0 { // found
			return i
		}
	}
	return -1 // not found
}
