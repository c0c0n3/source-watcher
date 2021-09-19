package nbic

import (
	"net/url"

	"github.com/fluxcd/source-watcher/osmops/util"
)

// Connection holds the data needed to establish a network connection with
// the OSM NBI.
type Connection struct {
	Address util.HostAndPort
	Secure  bool
}

func (b Connection) buildUrl(path string) (*url.URL, error) {
	return b.Address.BuildHttpUrl(b.Secure, path)
}

// Tokens returns the URL to the NBI tokens endpoint.
func (b Connection) Tokens() (*url.URL, error) {
	return b.buildUrl("/osm/admin/v1/tokens")
}
