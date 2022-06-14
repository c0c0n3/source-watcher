package nbic

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	u "github.com/fluxcd/source-watcher/osmops/util/http"
)

func stringReader(data string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(data))
}

type requestReply struct {
	req *http.Request
	res *http.Response
}

type mockNbi struct {
	handlers  map[string]u.ReqSender
	exchanges []requestReply
}

func newMockNbi() *mockNbi {
	mock := &mockNbi{
		handlers:  map[string]u.ReqSender{},
		exchanges: []requestReply{},
	}

	mock.handlers[handlerKey("POST", "/osm/admin/v1/tokens")] = tokenHandler
	mock.handlers[handlerKey("GET", "/osm/nsd/v1/ns_descriptors")] = nsDescHandler
	mock.handlers[handlerKey("GET", "/osm/admin/v1/vim_accounts")] = vimAccHandler
	mock.handlers[handlerKey("GET", "/osm/nslcm/v1/ns_instances_content")] = nsInstContentHandler
	mock.handlers[handlerKey("POST", "/osm/nslcm/v1/ns_instances_content")] = nsInstContentHandler
	mock.handlers[handlerKey("POST",
		"/osm/nslcm/v1/ns_instances/0335c32c-d28c-4d79-9b94-0ffa36326932/action")] = nsInstActionHandler
	mock.handlers[handlerKey("POST",
		"/osm/vnfpkgm/v1/vnf_packages_content")] = createPkgHandler
	mock.handlers[handlerKey("PUT",
		"/osm/vnfpkgm/v1/vnf_packages_content/my_knf")] = updatePkgHandler
	mock.handlers[handlerKey("POST",
		"/osm/nsd/v1/ns_descriptors_content")] = createPkgHandler
	mock.handlers[handlerKey("PUT",
		"/osm/nsd/v1/ns_descriptors_content/my_ns")] = updatePkgHandler

	return mock
}

func handlerKey(method string, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}

func (s *mockNbi) exchange(req *http.Request) (*http.Response, error) {
	key := handlerKey(req.Method, req.URL.Path)
	handle, ok := s.handlers[key]
	if !ok {
		return &http.Response{StatusCode: http.StatusInternalServerError},
			fmt.Errorf("no handler for request: %s", key)
	}

	res, err := handle(req)
	rr := requestReply{req: req, res: res}
	s.exchanges = append(s.exchanges, rr)

	return res, err
}

func tokenHandler(req *http.Request) (*http.Response, error) {
	reqCreds := UserCredentials{}
	json.NewDecoder(req.Body).Decode(&reqCreds)
	if reqCreds.Password != usrCreds.Password {
		return &http.Response{StatusCode: http.StatusUnauthorized}, nil
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       stringReader(validNbiTokenPayload),
	}, nil
}

func nsDescHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       stringReader(nsDescriptors),
	}, nil
}

func vimAccHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       stringReader(vimAccounts),
	}, nil
}

func nsInstContentHandler(req *http.Request) (*http.Response, error) {
	if req.Method == "GET" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       stringReader(nsInstancesContent),
		}, nil
	}

	// POST
	return &http.Response{StatusCode: http.StatusCreated}, nil
}

func nsInstActionHandler(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusAccepted}, nil
}

func createPkgHandler(req *http.Request) (*http.Response, error) {
	return echoReceivedStatusCode("POST", req)
}

func updatePkgHandler(req *http.Request) (*http.Response, error) {
	return echoReceivedStatusCode("PUT", req)
}

type StatusEcho struct {
	Code int
}

func echoReceivedStatusCode(method string, req *http.Request) (*http.Response,
	error) {
	if req.Method != method {
		return &http.Response{StatusCode: http.StatusMethodNotAllowed}, nil
	}

	inStatus := StatusEcho{}
	dec := gob.NewDecoder(req.Body)
	if err := dec.Decode(&inStatus); err != nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, nil
	}

	return &http.Response{
		StatusCode: int(inStatus.Code),
	}, nil
}
