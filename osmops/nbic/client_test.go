package nbic

import (
	"testing"
)

func TestNewNbicErrorOnNilTransport(t *testing.T) {
	if client, err := New(newConn(), usrCreds, nil); err == nil {
		t.Errorf("want: error; got: %+v", client)
	}
}

func TestNewNbicWithDefaultTransport(t *testing.T) {
	client, err := New(newConn(), usrCreds)
	if err != nil {
		t.Errorf("want: client; got: %v", err)
	}
	if client.transport == nil {
		t.Errorf("want: transport; got: nil")
	}
}
