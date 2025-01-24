package utils

import "testing"

func Test_generate2FA(t *testing.T) {
	username := "test"
	fa, s, err := Generate2FA(username)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("fa secret: %s, url: %s", fa, s)
}
