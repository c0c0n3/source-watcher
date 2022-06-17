package tgz

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"io/ioutil"
	"os"
	"testing"

	"github.com/fluxcd/source-watcher/osmops/util/file"
)

const ArchiveTestDirName = "archive_test_dir"

func findTestDataDir() file.AbsPath {
	_, thisFileName, _, _ := runtime.Caller(1)
	enclosingDir := filepath.Dir(thisFileName)
	testDataDir := filepath.Join(enclosingDir, ArchiveTestDirName)
	p, _ := file.ParseAbsPath(testDataDir)

	return p
}

func withTempDir(t *testing.T, do func(string)) {
	if tempDir, err := ioutil.TempDir("", "tgz-test"); err != nil {
		t.Errorf("couldn't create temp dir: %v", err)
	} else {
		defer os.RemoveAll(tempDir)
		defer os.Chmod(tempDir, 0700) // make sure you can remove it
		do(tempDir)
	}
}

func TestTgzThenExtract(t *testing.T) {
	withTempDir(t, func(tempDirPath string) {
		sourceDir := findTestDataDir()
		tarball, _ := file.ParseAbsPath(path.Join(tempDirPath, "test.tgz"))
		extractedDir := path.Join(tempDirPath, ArchiveTestDirName)

		MakeTarball(sourceDir, tarball)
		ExtractTarball(tarball, tempDirPath)

		want, _ := file.ListPaths(sourceDir.Value())
		got, _ := file.ListPaths(extractedDir)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v; got: %v", want, got)
		}
	})
}

func checkFileContent(pathname string) error {
	name := path.Base(pathname)
	content, err := ioutil.ReadFile(pathname)
	if err != nil {
		return err
	}
	text := string(content)
	if name != text {
		return fmt.Errorf("path = %s; content = %s", pathname, text)
	}
	return nil
}

func TestTgzThenExtractContent(t *testing.T) {
	withTempDir(t, func(tempDirPath string) {
		sourceDir := findTestDataDir()
		tarball, _ := file.ParseAbsPath(path.Join(tempDirPath, "test.tgz"))
		extractedDir, _ := file.ParseAbsPath(
			path.Join(tempDirPath, ArchiveTestDirName))

		MakeTarball(sourceDir, tarball)
		ExtractTarball(tarball, tempDirPath)

		scanner := file.NewTreeScanner(extractedDir)
		es := scanner.Visit(func(node file.TreeNode) error {
			if !node.FsMeta.IsDir() {
				return checkFileContent(node.NodePath.Value())
			}
			return nil
		})
		if len(es) > 0 {
			t.Errorf("want no errors; got: %v", es)
		}
	})
}
