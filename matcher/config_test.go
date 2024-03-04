package matcher

import (
	"testing"
)

func TestConfig(t *testing.T) {
	var myconfig, err = LoadConfig("../etc/data-matcher.yaml")
	if err != nil {
		t.Fatalf("loadConfig error %s", err)
	}
	t.Logf("myconfig: %s", myconfig.Dump())
}
