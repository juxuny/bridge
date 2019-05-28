package bridge

import "testing"

func TestTokenManager_Load(t *testing.T) {
	tokenManager := NewTokenManager()
	e := tokenManager.Load("token.conf")
	if e != nil {
		t.Error(e)
	}
	for _, item := range tokenManager.AuthToken {
		t.Log(item.Token, item.Key)
	}
}
