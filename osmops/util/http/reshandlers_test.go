package http

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type TestData struct {
	X int         `json:"x"`
	Y interface{} `json:"y"`
}

func stringReader(data string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(data))
}

func TestJsonReaderErrorOnNilResponse(t *testing.T) {
	target := TestData{}
	reader := ReadJsonResponse(&target)
	if err := reader.Handle(nil); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestJsonReaderErrorOnNilTarget(t *testing.T) {
	reader := ReadJsonResponse(nil)
	if err := reader.Handle(&http.Response{}); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestJsonReaderErrorOnUnexpectedResponseCode(t *testing.T) {
	target := TestData{}
	reader := ReadJsonResponse(&target) // default code: 200
	if err := reader.Handle(&http.Response{StatusCode: 400}); err == nil {
		t.Errorf("want: error; got: nil")
	}

	reader = ReadJsonResponse(&target, 201, 401) // no default code of 200
	if err := reader.Handle(&http.Response{StatusCode: 400}); err == nil {
		t.Errorf("want: error; got: nil")
	}
}

func TestJsonReaderGetData(t *testing.T) {
	target := TestData{}
	response := &http.Response{
		StatusCode: 200,
		Body:       stringReader(`{"x": 1, "y": {"z": 2}}`),
	}
	send := func(req *http.Request) (*http.Response, error) {
		return response, nil
	}
	res, err := Request(GET).
		SetHandler(ReadJsonResponse(&target)).
		RunWith(send)

	if err != nil {
		t.Errorf("want: deserialized JSON; got: %v", err)
	}
	if res != response {
		t.Errorf("want: %v; got: %v", response, res)
	}
	if target.X != 1.0 {
		t.Errorf("want: deserialized JSON; got: %+v", target)
	}
	if y, ok := target.Y.(map[string]interface{}); !ok {
		t.Errorf("want: deserialized JSON; got: %+v", target)
	} else {
		if y["z"] != 2.0 {
			t.Errorf("want: deserialized JSON; got: %+v", target)
		}
	}
}

func TestExpectSuccess(t *testing.T) {
	response := &http.Response{}
	send := func(req *http.Request) (*http.Response, error) {
		return response, nil
	}
	for code := 200; code < 300; code++ {
		response.StatusCode = code
		_, err := Request(GET).
			SetHandler(ExpectSuccess()).
			RunWith(send)
		if err != nil {
			t.Errorf("want: success; got: %v", err)
		}
	}
	for _, code := range []int{100, 199, 300, 400, 500} {
		response.StatusCode = code
		_, err := Request(GET).
			SetHandler(ExpectSuccess()).
			RunWith(send)
		if err == nil {
			t.Errorf("[%d] want: error; got: nil", code)
		}
	}
}

func TestExpectStatusCodeOneOf(t *testing.T) {
	response := &http.Response{}
	send := func(req *http.Request) (*http.Response, error) {
		return response, nil
	}
	want := []int{200, 201, 404}
	for _, code := range want {
		response.StatusCode = code
		_, err := Request(GET).
			SetHandler(ExpectStatusCodeOneOf(want...)).
			RunWith(send)
		if err != nil {
			t.Errorf("want: success; got: %v", err)
		}
	}
	for _, code := range []int{100, 199, 300, 400, 500} {
		response.StatusCode = code
		_, err := Request(GET).
			SetHandler(ExpectStatusCodeOneOf(want...)).
			RunWith(send)
		if err == nil {
			t.Errorf("[%d] want: error; got: nil", code)
		}
	}
}

func TestExpectStatusCodeNone(t *testing.T) {
	response := &http.Response{StatusCode: 200}
	send := func(req *http.Request) (*http.Response, error) {
		return response, nil
	}
	_, err := Request(GET).
		SetHandler(ExpectStatusCodeOneOf()).
		RunWith(send)
	if err == nil {
		t.Errorf("want: error; got: nil")
	}
}
