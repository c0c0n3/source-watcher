// Client to interact with OSM north-bound interface (NBI).
package nbic

import (
	"crypto/tls"
	"net/http"
	"time"

	. "github.com/fluxcd/source-watcher/osmops/util/http"
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
