package nbic

import (
	"testing"
)

// expired on Wed Sep 08 2021 18:52:11 GMT+0000
var expiredTokenPayload = `{
	"issued_at": 1631123531.1251214,
	"expires": 1631127131.1251214,
	"_id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"project_id": "fada443a-905c-4241-8a33-4dcdbdac55e7",
	"project_name": "admin",
	"username": "admin",
	"user_id": "5c6f2d64-9c23-4718-806a-c74c3fc3c98f",
	"admin": true,
	"roles": [{
		"name": "system_admin",
		"id": "cb545e44-cd2b-4c0b-93aa-7e2cee79afc3"
	}]
}`

// expires on Sat May 17 2053 20:38:51 GMT+0000
var validTokenPayload = `{
	"issued_at": 2631127131.1251214,
	"expires": 2631127131.1251214,
	"_id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"project_id": "fada443a-905c-4241-8a33-4dcdbdac55e7",
	"project_name": "admin",
	"username": "admin",
	"user_id": "5c6f2d64-9c23-4718-806a-c74c3fc3c98f",
	"admin": true,
	"roles": [{
		"name": "system_admin",
		"id": "cb545e44-cd2b-4c0b-93aa-7e2cee79afc3"
	}]
}`

func TestExpiredToken(t *testing.T) {
	token, err := FromNbiTokenPayload([]byte(expiredTokenPayload))
	if err != nil {
		t.Errorf("want: token; got: %v", err)
	}

	wantData := "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2"
	if token.String() != wantData {
		t.Errorf("want: %s; got: %s", wantData, token.String())
	}
	if !token.HasExpired() {
		t.Errorf("want: expired; got: still valid")
	}
}

func TestValidToken(t *testing.T) {
	token, err := FromNbiTokenPayload([]byte(validTokenPayload))
	if err != nil {
		t.Errorf("want: token; got: %v", err)
	}

	wantData := "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2"
	if token.String() != wantData {
		t.Errorf("want: %s; got: %s", wantData, token.String())
	}
	if token.HasExpired() {
		t.Errorf("want: still valid; got: expired")
	}
}

func TestEmptyTokenPayload(t *testing.T) {
	payload := `{"x": 1}`
	token, err := FromNbiTokenPayload([]byte(payload))
	if err != nil {
		t.Errorf("want: token; got: %v", err)
	}

	wantData := ""
	if token.String() != wantData {
		t.Errorf("want: %s; got: %s", wantData, token.String())
	}
	if !token.HasExpired() {
		t.Errorf("want: expired; got: still valid")
	}
}

func TestTokenPayloadDeserializationError(t *testing.T) {
	payload := `["expecting", "an object", "not an array!"]`
	if _, err := FromNbiTokenPayload([]byte(payload)); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestNullTokenPayloadDeserializationError(t *testing.T) {
	var payload []byte = nil
	if _, err := FromNbiTokenPayload(payload); err == nil {
		t.Errorf("want: error; got: nil")
	}
}
