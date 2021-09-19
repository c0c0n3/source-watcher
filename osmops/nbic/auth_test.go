package nbic

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/fluxcd/source-watcher/osmops/util"
)

// expired on Wed Sep 08 2021 18:52:11 GMT+0000
var expiredNbiTokenPayload = `{
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
var validNbiTokenPayload = `{
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

func stringReader(data string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(data))
}

var usrCreds = UserCredentials{
	Username: "admin", Password: "admin", Project: "admin",
}

func sameCreds(expected UserCredentials, req *http.Request) bool {
	got := UserCredentials{}
	json.NewDecoder(req.Body).Decode(&got)
	return reflect.DeepEqual(expected, got)
}

func newConn() Connection {
	address, _ := util.ParseHostAndPort("localhost:8080")
	return Connection{Address: *address}
}

type mockTransport struct {
	received    *http.Request
	replyWith   *http.Response
	timesCalled int
}

func (m *mockTransport) send(req *http.Request) (*http.Response, error) {
	m.received = req
	m.timesCalled += 1
	return m.replyWith, nil
}

func TestGetExpiredToken(t *testing.T) {
	mock := &mockTransport{
		replyWith: &http.Response{
			StatusCode: http.StatusOK,
			Body:       stringReader(expiredNbiTokenPayload),
		},
	}
	mngr, _ := NewAuthz(newConn(), usrCreds, mock.send)

	if _, err := mngr.GetAccessToken(); err == nil {
		t.Errorf("want: error; got: nil")
	}
	if mock.timesCalled != 1 {
		t.Errorf("want: 1; got: %d", mock.timesCalled)
	}
	if !sameCreds(usrCreds, mock.received) {
		t.Errorf("want: same usrCreds; got: different")
	}
}

func TestGetValidToken(t *testing.T) {
	mock := &mockTransport{
		replyWith: &http.Response{
			StatusCode: http.StatusOK,
			Body:       stringReader(validNbiTokenPayload),
		},
	}
	mngr, _ := NewAuthz(newConn(), usrCreds, mock.send)
	token, err := mngr.GetAccessToken()

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

	if mock.timesCalled != 1 {
		t.Errorf("want: 1; got: %d", mock.timesCalled)
	}
	if !sameCreds(usrCreds, mock.received) {
		t.Errorf("want: same usrCreds; got: different")
	}
}

func TestGetTokenPayloadWithNoTokenFields(t *testing.T) {
	mock := &mockTransport{
		replyWith: &http.Response{
			StatusCode: http.StatusOK,
			Body:       stringReader(`{"x": 1}`),
		},
	}
	mngr, _ := NewAuthz(newConn(), usrCreds, mock.send)
	if _, err := mngr.GetAccessToken(); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestGetTokenPayloadDeserializationError(t *testing.T) {
	mock := &mockTransport{
		replyWith: &http.Response{
			StatusCode: http.StatusOK,
			Body:       stringReader(`["expecting", "an object", "not an array!"]`),
		},
	}
	mngr, _ := NewAuthz(newConn(), usrCreds, mock.send)
	if _, err := mngr.GetAccessToken(); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestNewAuthzErrorOnNilTransport(t *testing.T) {
	if _, err := NewAuthz(Connection{}, UserCredentials{}, nil); err == nil {
		t.Errorf("want: error; got: nil")
	}
}
