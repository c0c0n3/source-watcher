// Client to interact with OSM north-bound interface (NBI).
package nbic

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	. "github.com/fluxcd/source-watcher/osmops/util/http"
	"github.com/fluxcd/source-watcher/osmops/util/http/sec"
)

const REQUEST_TIMEOUT_SECONDS = 600

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // (1)
			},
		},
		Timeout: time.Second * REQUEST_TIMEOUT_SECONDS, // (2)
	}
	// NOTE.
	// 1. Man-in-the-middle attacks. OSM client doesn't validate the server
	// cert, so we do the same. But this is a huge security loophole since it
	// opens the door to man-in-the-middle attacks.
	// 2. Request timeout. Always specify it, see
	// - https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
}

type Session struct {
	conn      Connection
	creds     UserCredentials
	transport ReqSender
	authz     *sec.TokenManager
	nsdMap    nsDescMap
	vimAccMap vimAccountMap
	nsInstMap nsInstanceMap
}

func New(conn Connection, creds UserCredentials, transport ...ReqSender) (
	*Session, error) {
	httpc := newHttpClient()

	agent := httpc.Do
	if len(transport) > 0 {
		agent = transport[0]
	}

	authz, err := NewAuthz(conn, creds, agent)
	if err != nil {
		return nil, err
	}

	return &Session{
		conn:      conn,
		creds:     creds,
		transport: agent,
		authz:     authz,
	}, nil
}

func (c *Session) NbiAccessToken() ReqBuilder {
	provider := func() (string, error) {
		if token, err := c.authz.GetAccessToken(); err != nil {
			return "", err
		} else {
			return token.String(), nil
		}
	}
	return BearerToken(provider)
}

func (c *Session) getJson(endpoint *url.URL, data interface{}) (
	*http.Response, error) {
	return Request(
		GET, At(endpoint),
		c.NbiAccessToken(),
		Accept(MediaType.JSON),
	).
		SetHandler(ExpectSuccess(), ReadJsonResponse(data)).
		RunWith(c.transport)
}

func (c *Session) postJson(endpoint *url.URL, inData interface{},
	outData ...interface{}) (*http.Response, error) {
	req := Request(
		POST, At(endpoint),
		c.NbiAccessToken(),
		Accept(MediaType.JSON),
		Content(MediaType.YAML), // same as what OSM client does
		JsonBody(inData),
	)
	if len(outData) > 0 {
		req.SetHandler(ExpectSuccess(), ReadJsonResponse(outData[0]))
	}
	return req.RunWith(c.transport)
}
