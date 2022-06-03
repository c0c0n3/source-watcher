package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/fluxcd/source-watcher/osmops/util/file"
)

// Writer writes data to a compressed tar stream.
// The tar format is PAX and the compression is gzip.
// You create a Writer with a sink stream where the compressed tar data
// gets written.
//
// Example. Archiving all the files in "some/dir" and its sub-directories.
//
//     sourceDir, _ := file.ParseAbsPath("some/dir")
//     sink, _ := os.OpenFile("my.tgz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
//     archiveBaseDirName := "my-root"
//
//     scanner := file.NewTreeScanner(sourceDir)
//     writer, _ := NewWriter(archiveBaseDirName, sink, gzip.BestCompression)
//
//     defer writer.Close()
//     scanner.Visit(writer.Visitor())
//
// If "some/dir/d1/f1" is a file, then it'll be archived at "my-root/d1/f1".
// Use an empty string for archive base directory name if you don't want to
// prefix archived paths---e.g. the above file would be archived at "d1/f1".
// The above example uses a file.Visitor to collect files from a directory,
// but you can also call directly the AddEntry and AddFile methods for finer
// control over what gets written to the archive. Also, there's a couple of
// convenience functions to archive directory contents to a stream or a file,
// see: WriteFileArchive and MakeTarball.
type Writer interface {
	// AddEntry writes the given content to the archive at the specified
	// path, relative to the archive base directory.
	AddEntry(archivePath string, content io.Reader) error
	// AddFile writes the given file to the archive at the specified path,
	// relative to the archive base directory.
	AddFile(archivePath string, filePath string, fi os.FileInfo) error
	// Visitor returns a function you can use with a file.TreeScanner to
	// collect all the files in a directory (including sub-directories)
	// and add them to the archive.
	Visitor() file.Visitor
	// Close finalises the writing to the archive and closes the underlying
	// sink stream.
	Close()
}

type tarball struct {
	baseDirName      string
	contentStream    *tar.Writer
	compressedStream *gzip.Writer
	sink             io.WriteCloser
}

func NewWriter(archiveBaseDirName string, sink io.WriteCloser,
	compressionLevel int) (Writer, error) {
	if sink == nil {
		return nil, fmt.Errorf("nil sink")
	}
	gzipStream, err := gzip.NewWriterLevel(sink, compressionLevel)
	if err != nil {
		return nil, err
	}
	tarStream := tar.NewWriter(gzipStream)

	return &tarball{
		baseDirName:      archiveBaseDirName,
		contentStream:    tarStream,
		compressedStream: gzipStream,
		sink:             sink,
	}, nil
}

func (t *tarball) Close() {
	t.contentStream.Close()
	t.compressedStream.Close()
	t.sink.Close()
}

func (t *tarball) writeHeader(archivePath string, hdr *tar.Header) error {
	hdr.Name = path.Join(t.baseDirName, archivePath)
	hdr.Typeflag = tar.TypeReg
	hdr.Format = tar.FormatPAX

	return t.contentStream.WriteHeader(hdr)
}

func (t *tarball) AddEntry(archivePath string, content io.Reader) error {
	header := &tar.Header{Mode: int64(0644)}
	if err := t.writeHeader(archivePath, header); err != nil {
		return err
	}
	if _, err := io.Copy(t.contentStream, content); err != nil {
		return err
	}
	return nil
}

func (t *tarball) AddFile(archivePath, filePath string, fi os.FileInfo) error {
	var err error

	if !fi.Mode().IsRegular() {
		return nil
	}

	header, err := tar.FileInfoHeader(fi, fi.Name())
	if err != nil {
		return err
	}
	err = t.writeHeader(archivePath, header)
	if err != nil {
		return err
	}

	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(t.contentStream, fd)
	fd.Close()

	return err
}

func (t *tarball) Visitor() file.Visitor {
	return func(node file.TreeNode) error {
		return t.AddFile(node.RelPath, node.NodePath.Value(), node.FsMeta)
	}
}
