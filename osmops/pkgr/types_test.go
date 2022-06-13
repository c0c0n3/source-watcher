package pkgr

import (
	"reflect"
	"testing"

	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
)

func TestFileHashFailedLookup(t *testing.T) {
	srcDir, _ := file.ParseAbsPath("no/where")
	pkgSrc := newPkgSrc(srcDir)
	got := pkgSrc.FileHash("not/there")
	if got != "" {
		t.Errorf("want: empty; got: %s", got)
	}
}

func TestReadAllData(t *testing.T) {
	srcDir, _ := file.ParseAbsPath("no/where")
	src := newPkgSrc(srcDir)
	data := bytez.NewBuffer()
	data.Write([]byte{1, 2, 3})
	pkg := makePackage(src, data)

	want := []byte{1, 2, 3}
	got := pkg.ReadAllData()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}
