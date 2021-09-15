package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	u "github.com/fluxcd/source-watcher/osmops/util"
)

func TestSimpleGetRequest(t *testing.T) {
	hp, _ := u.ParseHostAndPort("x:80")
	url, _ := hp.Http("/a/b")
	req, err := BuildRequest(
		GET, At(url),
	)

	if err != nil {
		t.Fatalf("want request, but got error: %v", err)
	}

	wantMethod := "GET"
	if req.Method != wantMethod {
		t.Errorf("want: %s; got: %s", wantMethod, req.Method)
	}

	wantUrl := "http://x:80/a/b"
	if req.URL.String() != wantUrl {
		t.Errorf("want: %s; got: %s", wantUrl, req.URL.String())
	}

	wantHost := "x:80"
	if req.Host != wantHost {
		t.Errorf("want: %s; got: %s", wantHost, req.Host)
	}
}

func TestSimplePostRequest(t *testing.T) {
	hp, _ := u.ParseHostAndPort("x:80")
	url, _ := hp.Http("/a/b")
	content := []byte("42")
	req, err := BuildRequest(
		POST, At(url),
		Body(content),
	)

	if err != nil {
		t.Fatalf("want request, but got error: %v", err)
	}

	wantMethod := "POST"
	if req.Method != wantMethod {
		t.Errorf("want: %s; got: %s", wantMethod, req.Method)
	}

	wantUrl := "http://x:80/a/b"
	if req.URL.String() != wantUrl {
		t.Errorf("want: %s; got: %s", wantUrl, req.URL.String())
	}

	wantHost := "x:80"
	if req.Host != wantHost {
		t.Errorf("want: %s; got: %s", wantHost, req.Host)
	}

	wantContentLength := int64(2)
	if req.ContentLength != wantContentLength {
		t.Errorf("want: %d; got: %d", wantContentLength, req.ContentLength)
	}

	gotBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("want: %v; got: %v", content, err)
	}
	if string(gotBody) != string(content) {
		t.Errorf("want: %v; got: %v", content, gotBody)
	}

	gotBodyReader, err := req.GetBody()
	if err != nil {
		t.Errorf("want %v; got: %v", content, err)
	}
	gotBody, err = ioutil.ReadAll(gotBodyReader)
	if err != nil {
		t.Errorf("want: %v; got: %v", content, err)
	}
	if string(gotBody) != string(content) {
		t.Errorf("want: %v; got: %v", content, gotBody)
	}
}

func TestEmptyBody(t *testing.T) {
	content := []byte("")
	req, err := BuildRequest(
		Body(content),
	)

	if err != nil {
		t.Fatalf("want request, but got error: %v", err)
	}

	wantContentLength := int64(0)
	if req.ContentLength != wantContentLength {
		t.Errorf("want: %d; got: %d", wantContentLength, req.ContentLength)
	}

	gotBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("want empty body; got: %v", err)
	}
	if string(gotBody) != "" {
		t.Errorf("want empty body; got: %v", gotBody)
	}

	gotBodyReader, err := req.GetBody()
	if err != nil {
		t.Errorf("want empty body; got: %v", err)
	}
	gotBody, err = ioutil.ReadAll(gotBodyReader)
	if err != nil {
		t.Errorf("want empty body; got: %v", err)
	}
	if string(gotBody) != "" {
		t.Errorf("want empty body; got: %v", gotBody)
	}
}

var acceptHeaderFixtures = []struct {
	in   []u.EnumIx
	want string
}{
	{in: []u.EnumIx{}, want: ""},
	{in: []u.EnumIx{MediaType.JSON}, want: "Accept: application/json\r\n"},
	{
		in:   []u.EnumIx{MediaType.JSON, MediaType.YAML},
		want: "Accept: application/json, application/yaml\r\n",
	},
}

func TestAcceptHeader(t *testing.T) {
	for k, d := range acceptHeaderFixtures {
		req, err := BuildRequest(
			Accept(d.in...),
		)
		if err != nil {
			t.Fatalf("[%d] want request, but got error: %v", k, err)
		}

		var buf bytes.Buffer
		if err := req.Header.Write(&buf); err != nil {
			t.Fatalf("[%d] want: %s; got: %v", k, d.want, err)
		}

		got := buf.String()
		if got != d.want {
			t.Errorf("[%d] want: %s; got: %s", k, d.want, got)
		}
	}
}

var contentTypeHeaderFixtures = []struct {
	in   u.EnumIx
	want string
}{
	{
		in:   MediaType.JSON,
		want: "Content-Type: application/json\r\n",
	},
	{
		in:   MediaType.YAML,
		want: "Content-Type: application/yaml\r\n",
	},
}

func TestContentTypeHeader(t *testing.T) {
	for k, d := range contentTypeHeaderFixtures {
		req, err := BuildRequest(
			Content(d.in),
		)
		if err != nil {
			t.Fatalf("[%d] want request, but got error: %v", k, err)
		}

		var buf bytes.Buffer
		if err := req.Header.Write(&buf); err != nil {
			t.Fatalf("[%d] want: %s; got: %v", k, d.want, err)
		}

		got := buf.String()
		if got != d.want {
			t.Errorf("[%d] want: %s; got: %s", k, d.want, got)
		}
	}
}

func TestBearerTokenHeader(t *testing.T) {
	tokenProvider := func() (string, error) { return "token", nil }
	req, err := BuildRequest(
		BearerToken(tokenProvider),
	)

	if err != nil {
		t.Fatalf("want request, but got error: %v", err)
	}

	want := "Authorization: Bearer token\r\n"

	var buf bytes.Buffer
	if err := req.Header.Write(&buf); err != nil {
		t.Fatalf("want: %s; got: %v", want, err)
	}

	got := buf.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestBearerTokenHeaderFail(t *testing.T) {
	tokenProvider := func() (string, error) { return "", errors.New("ouch!") }
	req, err := BuildRequest(
		BearerToken(tokenProvider),
	)

	if err == nil {
		t.Fatalf("want error, but got request: %v", req)
	}
}
