package nbic

import (
	"encoding/gob"
	"testing"

	"github.com/fluxcd/source-watcher/osmops/pkgr"
	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
)

func writeMockTgzStream(t *testing.T, httpStatusCodeToEcho int) *bytez.Buffer {
	data := StatusEcho{Code: httpStatusCodeToEcho}
	buf := bytez.NewBuffer()
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(data); err != nil {
		t.Fatalf("couldn't encode status echo: %v", err)
	}
	return buf
}

func runCreatePackageTest(t *testing.T, pkgId string) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)

	pkg := &pkgr.Package{
		Name:   pkgId,
		Source: nil,
		Data:   writeMockTgzStream(t, 201),
		Hash:   "h1",
	}
	handler := &pkgHandler{
		session: nbic,
		pack: func(src file.AbsPath) (*pkgr.Package, error) {
			return pkg, nil
		},
	}

	if err := handler.process(); err != nil {
		t.Errorf("want: create package; got: %v", err)
	}
}

func runUpdatePackageTest(t *testing.T, pkgId string) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)

	pkg := &pkgr.Package{
		Name:   pkgId,
		Source: nil,
		Data:   writeMockTgzStream(t, 409),
		Hash:   "h1",
	}
	handler := &pkgHandler{
		session: nbic,
		pack: func(src file.AbsPath) (*pkgr.Package, error) {
			return pkg, nil
		},
	}

	if err := handler.process(); err != nil {
		t.Errorf("want: update package; got: %v", err)
	}
}

func TestCreateKnfPackage(t *testing.T) {
	runCreatePackageTest(t, "my_knf")
}

func TestCreateNsPackage(t *testing.T) {
	runCreatePackageTest(t, "my_ns")
}

func TestUpdateKnfPackage(t *testing.T) {
	runUpdatePackageTest(t, "my_knf")
}

func TestUpdateNsPackage(t *testing.T) {
	runUpdatePackageTest(t, "my_ns")
}
