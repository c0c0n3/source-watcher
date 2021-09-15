package http

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
)

func TestExchangeRequestBuilderFailure(t *testing.T) {
	res, err := Request(GET, At(nil)).RunWith(http.DefaultClient.Do)
	if res != nil {
		t.Errorf("want nil response; got: %v", res)
	}
	if err == nil {
		t.Errorf("want request build error; got: nil")
	}
}

func TestExchangeStopOnSendFailure(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, errors.New("ouch!")
	}
	res, err := Request(GET).RunWith(send)
	if res == nil {
		t.Errorf("want empty response; got: nil")
	}
	if err == nil {
		t.Errorf("want send error; got: nil")
	}
}

func TestExchangeErrorsOnNilResponse(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return nil, nil
	}
	res, err := Request(GET).RunWith(send)
	if res != nil {
		t.Errorf("want nil response; got: %v", res)
	}
	if err == nil {
		t.Errorf("want handle error; got: nil")
	}
}

func TestExchangeNoHandleIfNoHandlers(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, nil
	}
	res, err := Request(GET).RunWith(send)
	if res == nil {
		t.Errorf("want empty response; got: nil")
	}
	if err != nil {
		t.Errorf("want no error; got: %v", err)
	}
}

type EmptyBody struct {
	bytes.Buffer
	closed bool
}

func (e *EmptyBody) Close() error {
	e.closed = true
	return nil
}

func TestExchangeCloseBody(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return &http.Response{Body: &EmptyBody{}}, nil
	}
	res, err := Request(GET).RunWith(send)
	if res == nil {
		t.Errorf("want response; got: nil")
	}
	if err != nil {
		t.Errorf("want no error; got: %v", err)
	}

	body := (res.Body).(*EmptyBody)
	if !body.closed {
		t.Errorf("didn't close body stream on exit")
	}
}

type GrabStatusCode struct {
	code int
}

func (x *GrabStatusCode) Handle(response *http.Response) error {
	x.code = response.StatusCode
	return nil
}

type FailingHandler struct{}

func (x *FailingHandler) Handle(response *http.Response) error {
	return errors.New("ouch!")
}

func TestExchangeHandleResponseSuccessfully(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	}
	statusCodeGrabber := &GrabStatusCode{}
	res, err := Request(GET).
		SetHandler(statusCodeGrabber).
		RunWith(send)

	if res == nil {
		t.Errorf("want response; got: nil")
	}
	if err != nil {
		t.Errorf("want no error; got: %v", err)
	}
	if statusCodeGrabber.code != 200 {
		t.Errorf("want: 200; got: %d", statusCodeGrabber.code)
	}
}

func TestExchangeHandleResponseFailure(t *testing.T) {
	send := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	}
	statusCodeGrabber := &GrabStatusCode{}
	res, err := Request(GET).
		SetHandler(&FailingHandler{}, statusCodeGrabber).
		RunWith(send)

	if res == nil {
		t.Errorf("want response; got: nil")
	}
	if err == nil {
		t.Errorf("want error; got: nil")
	}
	if statusCodeGrabber.code == 200 {
		t.Errorf("want: don't run handlers following failed one")
	}
}

func TestRunWithNilReqSender(t *testing.T) {
	if _, err := Request().RunWith(nil); err == nil {
		t.Errorf("want error; got: nil")
	}
	HandleResponse(&http.Response{}, nil)
}
