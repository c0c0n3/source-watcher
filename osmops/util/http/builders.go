// Common ReqBuilder functions.

package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	u "github.com/fluxcd/source-watcher/osmops/util"
)

var GET ReqBuilder = func(request *http.Request) error {
	request.Method = "GET"
	return nil
}

var POST = func(request *http.Request) error {
	request.Method = "POST"
	return nil
}

func At(url *url.URL) ReqBuilder {
	return func(request *http.Request) error {
		if url == nil {
			return errors.New("nil URL")
		}
		request.URL = url
		request.Host = url.Host
		return nil
	}
}

var MediaType = struct {
	u.StrEnum
	JSON, YAML u.EnumIx
}{
	StrEnum: u.NewStrEnum("application/json", "application/yaml"),
	JSON:    0,
	YAML:    1,
}

func Content(mediaType u.EnumIx) ReqBuilder {
	return func(request *http.Request) error {
		request.Header.Set("Content-Type", MediaType.LabelOf(mediaType))
		return nil
	}
}

func Accept(mediaType ...u.EnumIx) ReqBuilder {
	return func(request *http.Request) error {
		ts := []string{}
		for _, mt := range mediaType {
			ts = append(ts, MediaType.LabelOf(mt))
		}
		if len(ts) > 0 {
			request.Header.Set("Accept", strings.Join(ts, ", "))
		}

		return nil
	}
	// TODO implement weights too? Not needed for OSM client.
}

func Authorization(value string) ReqBuilder {
	return func(request *http.Request) error {
		request.Header.Set("Authorization", value)
		return nil
	}
}

type BearerTokenProvider func() (string, error)

func BearerToken(acquireToken BearerTokenProvider) ReqBuilder {
	return func(request *http.Request) error {
		if token, err := acquireToken(); err != nil {
			return err
		} else {
			authValue := fmt.Sprintf("Bearer %s", token)
			return Authorization(authValue)(request)
		}
	}
}

func Body(content []byte) ReqBuilder {
	return func(request *http.Request) error {
		request.ContentLength = int64(len(content))

		if len(content) == 0 {
			// see code comments in Request.NewRequestWithContext about an
			// empty body and backward compat.
			request.Body = http.NoBody
			request.GetBody = func() (io.ReadCloser, error) {
				return http.NoBody, nil
			}
		} else {
			request.Body = io.NopCloser(bytes.NewBuffer(content))

			// the following code does the same as Request.NewRequestWithContext
			// so 307 and 308 redirects can replay the body.
			request.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(content)
				return io.NopCloser(r), nil
			}
		}

		return nil
	}
}

// TODO also implement streaming body? most of the standard libs aren't built
// w/ streaming in mind, so in practice you'll likely have the whole body in
// memory most of the time for common cases---e.g. JSON, YAML.

// TODO nil pointer checks. Mostly not implemented!! Catch all occurrences
// of slices, i/f, function args and return an error if nil gets passed in.
// Then write test cases for each. What a schlep!
