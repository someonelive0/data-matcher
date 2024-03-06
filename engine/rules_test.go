package engine

import "testing"

func TestRules(t *testing.T) {
	rulesConf, err := NewRulesConfig("../etc/rules.yaml.tpl")
	if err != nil {
		t.Fatalf("NewRulesConfig failed %s", err)
	}

	// t.Logf("%#v", rulesConf)

	regs := rulesConf.GetValueReg()
	t.Logf("%d, %#v", len(regs), regs)

	dicts := rulesConf.GetColDict()
	t.Logf("%d, %#v", len(dicts), dicts)
}
