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

func (b Connection) buildUrl(path string) *url.URL {
	if url, err := b.Address.BuildHttpUrl(b.Secure, path); err != nil {
		panic(err) // see note below
	} else {
		return url
	}
}

// NOTE. Panic on URL building.
// Ideally buildUrl should return (*url.URL, error) instead of panicing. But
// then it becomes a royal pain in the backside to write code that uses the
// URL functions below and testing for the URL build error case needs to
// happen at every calling site---e.g. if you call Tokens then ideally there
// should be a unit test to check what happens when Tokens returns an error.
// So we take a shortcut with the panic call. As long as we call all the
// URL functions below in our unit tests, we can be sure the panic won't
// happen at runtime.

// Tokens returns the URL to the NBI tokens endpoint.
func (b Connection) Tokens() *url.URL {
	return b.buildUrl("/osm/admin/v1/tokens")
}

// NsDescriptors returns the URL to the NBI NS descriptors endpoint.
func (b Connection) NsDescriptors() *url.URL {
	return b.buildUrl("/osm/nsd/v1/ns_descriptors")
}
