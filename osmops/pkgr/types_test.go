package pkgr

import (
	"testing"

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
