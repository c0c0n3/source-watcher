package pkgr

import (
	"compress/gzip"
	"io"

	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
	"github.com/fluxcd/source-watcher/osmops/util/tgz"
)

func Pack(source file.AbsPath) (*Package, error) {
	sink := bytez.NewBuffer()
	pkgSource := newPkgSrc(source)
	if err := writePackageData(pkgSource, sink); err != nil {
		return nil, err
	}
	return makePackage(pkgSource, sink), nil
}

func writePackageData(source *pkgSrc, sink io.WriteCloser) error {
	archiveBaseDirName := source.DirectoryName()
	writer, err := tgz.NewWriter(archiveBaseDirName, sink, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := collectPackageItems(source, writer); err != nil {
		return err
	}
	return addChecksumFile(source, writer)
}

func collectPackageItems(source *pkgSrc, writer tgz.Writer) error {
	scanner := file.NewTreeScanner(source.Directory())
	visitor := makeSourceVisitor(source, writer)
	if es := scanner.Visit(visitor); len(es) > 0 {
		return es[0]
	}
	return nil
}

func makeSourceVisitor(source *pkgSrc, writer tgz.Writer) file.Visitor {
	collectFile := writer.Visitor()
	return func(node file.TreeNode) error {
		if err := collectFile(node); err != nil {
			return err
		}
		return source.addFileHash(node)
	}
}

func addChecksumFile(source *pkgSrc, writer tgz.Writer) error {
	content := writeCheckSumFileContent(source)
	return writer.AddEntry(ChecksumFileName, content)
}
