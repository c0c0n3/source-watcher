package pkgr

import (
	"io"
	"path"

	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
)

type Package struct {
	Name   string
	Source file.AbsPath
	Data   io.ReadCloser
	Hash   string
}

func extractPackageName(src file.AbsPath) string {
	return path.Base(src.Value())
}

func makePackage(src file.AbsPath, data *bytez.Buffer) *Package {
	return &Package{
		Name:   extractPackageName(src),
		Source: src,
		Data:   data,
		Hash:   md5string(data.Bytes()),
	}
}

type pathToHash map[string]string
