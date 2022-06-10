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
	if err := writePackageData(source, sink); err != nil {
		return nil, err
	}
	return makePackage(source, sink), nil
}

func writePackageData(source file.AbsPath, sink io.WriteCloser) error {
	archiveBaseDirName := extractPackageName(source)
	writer, err := tgz.NewWriter(archiveBaseDirName, sink, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer writer.Close()

	m := make(pathToHash)
	if err := collectPackageItems(writer, m, source); err != nil {
		return err
	}
	err = addCheckSumFile(writer, m, archiveBaseDirName)

	return err
}

func collectPackageItems(writer tgz.Writer, m pathToHash, source file.AbsPath) error {
	scanner := file.NewTreeScanner(source)
	visitor := makeSourceVisitor(writer, m)
	if es := scanner.Visit(visitor); len(es) > 0 {
		return es[0]
	}
	return nil
}

func makeSourceVisitor(writer tgz.Writer, m pathToHash) file.Visitor {
	collectFile := writer.Visitor()
	return func(node file.TreeNode) error {
		if err := collectFile(node); err != nil {
			return err
		}
		if node.FsMeta.Mode().IsRegular() {
			if hash, err := computeCheckSum(node.NodePath); err != nil {
				return err
			} else {
				m[node.RelPath] = hash
			}
		}

		return nil
	}
}

func addCheckSumFile(writer tgz.Writer, m pathToHash, archiveBaseDirName string) error {
	content := writeCheckSumFileContent(archiveBaseDirName, m)
	if err := writer.AddEntry(ChecksumFileName, content); err != nil {
		return err
	}
	return nil
}
