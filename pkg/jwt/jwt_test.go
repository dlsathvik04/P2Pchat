package jwt_test

import (
	"testing"
	"time"

	"github.com/dlsathvik04/P2Pchat/pkg/jwt"
)

func TestSHA256JWTManager(t *testing.T) {
	testjwtm := jwt.NewJWTManager(time.Minute, "Secret", "onlytexttest")
	type testPayload struct {
		Userid   int
		Verified bool
	}
	tpl := testPayload{
		Userid:   8,
		Verified: true,
	}
	token := testjwtm.GenerateToken(tpl)

	var pl testPayload
	testjwtm.AuthorizeToken(token, &pl)

	if pl != tpl {
		t.Error("failed")
	}
}
