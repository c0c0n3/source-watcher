package nbic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fluxcd/source-watcher/osmops/pkgr"
	"github.com/fluxcd/source-watcher/osmops/util/file"
	. "github.com/fluxcd/source-watcher/osmops/util/http"
)

func (s *Session) CreateOrUpdatePackage(source file.AbsPath) error {
	reader, err := newPkgReader(source)
	if err != nil {
		return err
	}
	handler := pkgHandler{
		session: s,
		pkg:     reader,
	}
	return handler.process()
}

// pkgReader wraps Package to consolidate in one place all the assumptions
// this module makes about OSM packages in an OsmOps-managed repo.
// Specifically:
//
// - pkg name = pkg ID
// - VNF pkg => pgk name ends w/ "_knf"
// - NS pkg => pkg name ends w/ "_ns"
//
// None of the above needs to be true in general, but OsmOps relies on that
// at the moment to simplify the implementation. Eventually, we'll redo this
// properly, i.e. use a semantic approach (parse, interpret OSM files) rather
// than naming conventions and guesswork.
type pkgReader struct {
	pkg  *pkgr.Package
	data []byte
}

func newPkgReader(pkgSource file.AbsPath) (*pkgReader, error) {
	pkg, err := pkgr.Pack(pkgSource)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(pkg.Data)
	return &pkgReader{
		pkg:  pkg,
		data: data,
	}, err
}

func (r *pkgReader) Source() file.AbsPath {
	return r.pkg.Source.Directory()
}

func (r *pkgReader) Name() string {
	return r.pkg.Name
}

func (r *pkgReader) Id() string {
	return r.Name()
}

func (r *pkgReader) Data() []byte {
	return r.data
}

func (r *pkgReader) Hash() string {
	return r.pkg.Hash
}

func (r *pkgReader) IsNs() bool {
	return strings.HasSuffix(r.Name(), "_ns")
}

func (r *pkgReader) IsKnf() bool {
	return strings.HasSuffix(r.Name(), "_knf")
}

type pkgHandler struct {
	session *Session
	pkg     *pkgReader
}

func (h *pkgHandler) lookupCreateUrl() (*url.URL, error) {
	if h.pkg.IsKnf() {
		return h.session.conn.VnfPackagesContent(), nil
	}
	if h.pkg.IsNs() {
		return h.session.conn.NsPackagesContent(), nil
	}
	return nil, h.unsupportedPackageType()
}

func (h *pkgHandler) lookupUpdateUrl() (*url.URL, error) {
	if h.pkg.IsKnf() {
		return h.session.conn.VnfPackageContent(h.pkg.Id()), nil
	}
	if h.pkg.IsNs() {
		return h.session.conn.NsPackageContent(h.pkg.Id()), nil
	}
	return nil, h.unsupportedPackageType()
}

func (h *pkgHandler) unsupportedPackageType() error {
	return fmt.Errorf("unsupported package type: %v", h.pkg.Source())
}

func (h *pkgHandler) process() error {
	res, err := h.create()
	if err != nil {
		return err
	}

	if h.shouldUpdate(res) {
		_, err = h.update()
	}
	return err
}

func (h *pkgHandler) create() (*http.Response, error) {
	endpoint, err := h.lookupCreateUrl()
	if err != nil {
		return nil, err
	}
	return h.post(endpoint)
}

func (h *pkgHandler) shouldUpdate(res *http.Response) bool {
	return res.StatusCode == 409
}

func (h *pkgHandler) update() (*http.Response, error) {
	endpoint, err := h.lookupUpdateUrl()
	if err != nil {
		return nil, err
	}
	return h.put(endpoint)
}

func (h *pkgHandler) post(endpoint *url.URL) (*http.Response, error) {
	req := Request(
		POST, At(endpoint),
		h.session.NbiAccessToken(),
		Accept(MediaType.JSON),  // same as what OSM client does
		Content(MediaType.GZIP), // ditto
		ContentFilename(h.pkg),  // ditto
		ContentFileMd5(h.pkg),   // ditto
		Body(h.pkg.Data()),
	)
	req.SetHandler(ExpectStatusCodeOneOf(201, 409))
	return req.RunWith(h.session.transport)
}

func (h *pkgHandler) put(endpoint *url.URL) (*http.Response, error) {
	req := Request(
		PUT, At(endpoint),
		h.session.NbiAccessToken(),
		Accept(MediaType.JSON),  // same as what OSM client does
		Content(MediaType.GZIP), // ditto
		ContentFilename(h.pkg),  // ditto
		ContentFileMd5(h.pkg),   // ditto
		Body(h.pkg.Data()),
	)
	req.SetHandler(ExpectSuccess())
	return req.RunWith(h.session.transport)
}

func ContentFilename(pkg *pkgReader) ReqBuilder {
	name := fmt.Sprintf("%s.tar.gz", pkg.Name())
	return func(request *http.Request) error {
		request.Header.Set("Content-Filename", name)
		return nil
	}
}

func ContentFileMd5(pkg *pkgReader) ReqBuilder {
	return func(request *http.Request) error {
		request.Header.Set("Content-File-MD5", pkg.Hash())
		return nil
	}
}
