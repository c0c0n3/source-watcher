package tgz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fluxcd/source-watcher/osmops/util/file"
)

// ExtractTarball extracts the files in the given tarball to the specified
// directory, taking care of creating intermediate directories as needed.
//
// If you pass the empty string for destDirPath, ExtractTarball preserves
// the original archive paths, even if they're absolute. For example, if
// "/d/f" is the path of file f in the archive, ExtractTarball will try
// creating a directory "/d" if it doesn't exist and then put f in there.
// With an empty destDirPath, ExtractTarball resolves relative archive paths
// with respect to the current directory. For example, if "d/f" is the path
// of file f in the archive, ExtractTarball will try creating a directory
// "./d" if it doesn't exist and then put f in there.
//
// On the other hand, if you specify a destDirPath (either absolute or
// relative to the current directory), ExtractTarball recreates the directory
// structure of the archived files entirely in destDirPath by interpreting
// all archive paths (even absolute ones) relative to destDirPath. For example,
// if "/d/f" is the path of file f in the archive, ExtractTarball will try
// creating a directory "destDirPath/d" if it doesn't exist and then put f
// in there. The same happens to relative paths. For example, if "d/f" is
// the path of file f in the archive, ExtractTarball will try creating a
// directory "destDirPath/d" if it doesn't exist and then put f in there.
func ExtractTarball(tarballPath file.AbsPath, destDirPath string) error {
	source, err := os.Open(tarballPath.Value())
	if err != nil {
		return err
	}
	defer source.Close()

	deflate, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer deflate.Close()

	archive := tar.NewReader(deflate)
	return extractEntries(archive, destDirPath)
}

func extractEntries(archive *tar.Reader, destDirPath string) error {
	for {
		header, err := archive.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = extractEntry(archive, header, destDirPath)
		if err != nil {
			return err
		}
	}
}

func extractEntry(ar *tar.Reader, hdr *tar.Header, destDirPath string) error {
	targetPath := filepath.Join(destDirPath, hdr.Name)
	if err := ensureDirs(hdr, targetPath); err != nil {
		return err
	}

	fd, err := os.OpenFile(targetPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, hdr.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, ar)
	if err != nil {
		return err
	}
	return nil
}

func ensureDirs(hdr *tar.Header, targetPath string) error {
	if hdr.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, hdr.FileInfo().Mode())
	}
	enclosingDir := filepath.Dir(targetPath)
	return os.MkdirAll(enclosingDir, fs.ModePerm|fs.ModeDir)
}
