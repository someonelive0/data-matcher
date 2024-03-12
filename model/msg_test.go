package model

import (
	"encoding/json"
	"testing"
)

func TestMstHttp(t *testing.T) {
	s := `{}`

	msgHttp := &MsgHttp{}
	err := json.Unmarshal([]byte(s), msgHttp)
	if err != nil {
		t.Fatalf("TestMstHttp failed %s", err)
	}
	t.Logf("%d, %#v", len(s), *msgHttp)
}
