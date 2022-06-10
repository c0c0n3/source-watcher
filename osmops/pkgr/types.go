package pkgr

import (
	"io"
	"path"
	"sort"

	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
)

type Package struct {
	Name   string
	Source PackageSource
	Data   io.ReadCloser
	Hash   string
}

func makePackage(src PackageSource, data *bytez.Buffer) *Package {
	return &Package{
		Name:   src.DirectoryName(),
		Source: src,
		Data:   data,
		Hash:   md5string(data.Bytes()),
	}
}

type PackageSource interface {
	Directory() file.AbsPath
	DirectoryName() string
	SortedFilePaths() []string
	FileHash(filePath string) string
}

type pkgSrc struct {
	srcDir        file.AbsPath
	srcDirName    string
	pathToHashMap map[string]string
}

func newPkgSrc(srcDir file.AbsPath) *pkgSrc {
	return &pkgSrc{
		srcDir:        srcDir,
		srcDirName:    path.Base(srcDir.Value()),
		pathToHashMap: make(map[string]string),
	}
}

func (p *pkgSrc) Directory() file.AbsPath {
	return p.srcDir
}

func (p *pkgSrc) DirectoryName() string {
	return p.srcDirName
}

func (p *pkgSrc) SortedFilePaths() []string {
	keys := make([]string, 0, len(p.pathToHashMap))
	for k := range p.pathToHashMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (p *pkgSrc) FileHash(filePath string) string {
	if hash, ok := p.pathToHashMap[filePath]; ok {
		return hash
	}
	return ""
}

func (p *pkgSrc) addFileHash(node file.TreeNode) error {
	if !node.FsMeta.Mode().IsRegular() {
		return nil
	}
	hash, err := computeChecksum(node.NodePath)
	if err == nil {
		baseNamePlusPath := path.Join(p.srcDirName, node.RelPath)
		p.pathToHashMap[baseNamePlusPath] = hash
	}
	return err
}
