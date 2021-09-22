package nbic

import (
	"reflect"
	"testing"
)

func TestBuildNsInstanceMap(t *testing.T) {
	vs := []nsInstanceView{
		{Id: "1", Name: "a"}, {Id: "2", Name: "a"}, {Id: "3", Name: "b"},
	}
	nsMap := buildNsInstanceMap(vs)

	if got, ok := nsMap["a"]; !ok {
		t.Errorf("want: a; got: nil")
	} else {
		want := []string{"1", "2"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v; got: %v", want, got)
		}
	}

	if got, ok := nsMap["b"]; !ok {
		t.Errorf("want: b; got: nil")
	} else {
		want := []string{"3"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v; got: %v", want, got)
		}
	}

	if got, ok := nsMap["c"]; ok {
		t.Errorf("want: nil; got: %v", got)
	}
}

func TestLookupNsInstIdUseCachedData(t *testing.T) {
	nbic := &Session{
		nsInstMap: map[string][]string{"silly_ns": {"324567"}},
	}
	id, err := nbic.lookupNsInstanceId("silly_ns")
	if err != nil {
		t.Errorf("want: 324567; got: %v", err)
	}
	if id != "324567" {
		t.Errorf("want: 324567; got: %s", id)
	}
}

func TestLookupNsInstIdErrorOnMiss(t *testing.T) {
	nbic := &Session{
		nsInstMap: map[string][]string{"silly_ns": {"324567"}},
	}
	if _, err := nbic.lookupNsInstanceId("not there!"); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestLookupNsInstIdFetchDataFromServer(t *testing.T) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)

	wantId := "0335c32c-d28c-4d79-9b94-0ffa36326932"
	if gotId, err := nbic.lookupNsInstanceId("ldap"); err != nil {
		t.Errorf("want: %s; got: %v", wantId, err)
	} else {
		if gotId != wantId {
			t.Errorf("want: %s; got: %v", wantId, gotId)
		}
	}

	if len(nbi.exchanges) != 2 {
		t.Fatalf("want: 2; got: %d", len(nbi.exchanges))
	}
	rr1, rr2 := nbi.exchanges[0], nbi.exchanges[1]
	if rr1.req.URL.Path != urls.Tokens().Path {
		t.Errorf("want: %s; got: %s", urls.Tokens().Path, rr1.req.URL.Path)
	}
	if rr2.req.URL.Path != urls.NsInstancesContent().Path {
		t.Errorf("want: %s; got: %s", urls.NsInstancesContent().Path, rr2.req.URL.Path)
	}
}

func TestLookupNsInstIdFetchDataFromServerTokenError(t *testing.T) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, UserCredentials{}, nbi.exchange)

	if _, err := nbic.lookupNsInstanceId("ldap"); err == nil {
		t.Errorf("want: error; got: nil")
	}

	if len(nbi.exchanges) != 1 {
		t.Fatalf("want: 1; got: %d", len(nbi.exchanges))
	}
	rr1 := nbi.exchanges[0]
	if rr1.req.URL.Path != urls.Tokens().Path {
		t.Errorf("want: %s; got: %s", urls.Tokens().Path, rr1.req.URL.Path)
	}
}

func TestLookupNsInstIdFetchDataFromServerDupNameError(t *testing.T) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)

	if _, err := nbic.lookupNsInstanceId("dup-name"); err == nil {
		t.Errorf("want: error; got: nil")
	}

	if len(nbi.exchanges) != 2 {
		t.Fatalf("want: 2; got: %d", len(nbi.exchanges))
	}
	rr1, rr2 := nbi.exchanges[0], nbi.exchanges[1]
	if rr1.req.URL.Path != urls.Tokens().Path {
		t.Errorf("want: %s; got: %s", urls.Tokens().Path, rr1.req.URL.Path)
	}
	if rr2.req.URL.Path != urls.NsInstancesContent().Path {
		t.Errorf("want: %s; got: %s", urls.NsInstancesContent().Path, rr2.req.URL.Path)
	}
}
