package nbic

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fluxcd/source-watcher/osmops/pkgr"
	"github.com/fluxcd/source-watcher/osmops/util/file"
	. "github.com/fluxcd/source-watcher/osmops/util/http"
)

func (s *Session) CreateOrUpdatePackage(source file.AbsPath) error {
	handler := pkgHandler{s, source, pkgr.Pack}
	return handler.process()
}

type pkgHandler struct {
	session   *Session
	pkgSource file.AbsPath
	pack      func(file.AbsPath) (*pkgr.Package, error) // added for testability
}

// Assumptions:
// - pkg name = pkg ID
// - VNF pkg => pgk name ends w/ _knf
// - NS pkg => pkg name ends w/ _ns

func (h *pkgHandler) lookupCreateUrl(pkg *pkgr.Package) (*url.URL, error) {
	if strings.HasSuffix(pkg.Name, "_knf") {
		return h.session.conn.VnfPackagesContent(), nil
	}
	if strings.HasSuffix(pkg.Name, "_ns") {
		return h.session.conn.NsPackagesContent(), nil
	}
	return nil, fmt.Errorf("unsupported package type: %v", h.pkgSource)
}

func (h *pkgHandler) lookupUpdateUrl(pkg *pkgr.Package) (*url.URL, error) {
	if strings.HasSuffix(pkg.Name, "_knf") {
		return h.session.conn.VnfPackageContent(pkg.Name), nil
	}
	if strings.HasSuffix(pkg.Name, "_ns") {
		return h.session.conn.NsPackageContent(pkg.Name), nil
	}
	return nil, fmt.Errorf("unsupported package type: %v", h.pkgSource)
}

func (h *pkgHandler) process() error {
	pkg, err := h.pack(h.pkgSource)
	if err != nil {
		return err
	}

	if res, err := h.create(pkg); err == nil {
		if h.shouldUpdate(res) {
			_, err = h.update(pkg)
			return err
		}
	}
	return err
}

func (h *pkgHandler) create(pkg *pkgr.Package) (*http.Response, error) {
	endpoint, err := h.lookupCreateUrl(pkg)
	if err != nil {
		return nil, err
	}
	return h.post(endpoint, pkg)
}

func (h *pkgHandler) shouldUpdate(res *http.Response) bool {
	return res.StatusCode == 409
}

func (h *pkgHandler) update(pkg *pkgr.Package) (*http.Response, error) {
	endpoint, err := h.lookupUpdateUrl(pkg)
	if err != nil {
		return nil, err
	}
	return h.put(endpoint, pkg)
}

func (h *pkgHandler) post(endpoint *url.URL, pkg *pkgr.Package) (*http.Response,
	error) {
	req := Request(
		POST, At(endpoint),
		h.session.NbiAccessToken(),
		Accept(MediaType.JSON),  // same as what OSM client does
		Content(MediaType.GZIP), // ditto
		ContentFilename(pkg),    // ditto
		ContentFileMd5(pkg),     // ditto
		Body(pkg.ReadAllData()), // TODO stream instead?
	)
	req.SetHandler(ExpectStatusCodeOneOf(201, 409))
	return req.RunWith(h.session.transport)
}

func (h *pkgHandler) put(endpoint *url.URL, pkg *pkgr.Package) (*http.Response,
	error) {
	req := Request(
		PUT, At(endpoint),
		h.session.NbiAccessToken(),
		Accept(MediaType.JSON),  // same as what OSM client does
		Content(MediaType.GZIP), // ditto
		ContentFilename(pkg),    // ditto
		ContentFileMd5(pkg),     // ditto
		Body(pkg.ReadAllData()), // TODO stream instead?
	)
	req.SetHandler(ExpectSuccess())
	return req.RunWith(h.session.transport)
}

func ContentFilename(pkg *pkgr.Package) ReqBuilder {
	name := fmt.Sprintf("%s.tar.gz", pkg.Name)
	return func(request *http.Request) error {
		request.Header.Set("Content-Filename", name)
		return nil
	}
}

func ContentFileMd5(pkg *pkgr.Package) ReqBuilder {
	return func(request *http.Request) error {
		request.Header.Set("Content-File-MD5", pkg.Hash)
		return nil
	}
}
