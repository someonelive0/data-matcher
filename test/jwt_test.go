package test

import (
	"testing"

	"data-matcher/utils"
)

// TestJwt is test createToken and parseToken
func TestJwt(t *testing.T) {
	myJwt := utils.NewMyJwt(nil)
	tokenString, err := myJwt.CreateToken("testid")
	if err != nil {
		t.Fatalf("createToken failed: %s", err)
	}
	t.Logf("createToken: %s", tokenString)

	c, err := myJwt.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("parseToken failed: %s", err)
	}
	t.Logf("parseToken: %#v", *c)
	t.Logf("parseToken statand cleams: %#v", *c.StandardClaims)
}
