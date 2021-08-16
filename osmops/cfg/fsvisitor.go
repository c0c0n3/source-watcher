// Traversal of the repo target directory tree to process the content of any
// OSM GitOps files found in there.
//
package cfg

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	u "github.com/fluxcd/source-watcher/osmops/util"
)

// KduNsActionFile is the data passed to the OSM GitOps file visitor.
type KduNsActionFile struct {
	FilePath u.AbsPath
	Content  *KduNsAction
}

// KduNsActionProcessor is a file visitor that is given, in turn, the content
// of each OSM GitOps file found in the target directory.
type KduNsActionProcessor interface {
	// Do something with the current OSM GitOps file, possibly returning an
	// error if something goes wrong.
	process(file *KduNsActionFile) error
}

// VisitError wraps any error that happened while traversing the target
// directory with an additional path to indicate where the error happened.
type VisitError struct {
	AbsPath string
	Err     error
}

// Error implements the standard error interface.
func (e VisitError) Error() string {
	return fmt.Sprintf("%s: %v", e.AbsPath, e.Err)
}

// Unwrap implements Go's customary error unwrapping.
func (e VisitError) Unwrap() error { return e.Err }

// KduNsActionRepoScanner has methods to let visitors process OSM GitOps files
// found while traversing the target directory.
type KduNsActionRepoScanner struct {
	targetDir u.AbsPath
	fileExt   []u.NonEmptyStr
	parsePath func(string) (u.AbsPath, error) // (*)
	readFile  func(string) ([]byte, error)    // (*)

	// (*) added for testability, so we can sort of mock stuff
}

// NewKduNsActionRepoScanner instantiates a KduNsActionRepoScanner to
// traverse the target directory configured in the given Store.
func NewKduNsActionRepoScanner(store *Store) *KduNsActionRepoScanner {
	return &KduNsActionRepoScanner{
		targetDir: store.RepoTargetDirectory(),
		fileExt:   store.OpsFileExtensions(),
		parsePath: u.ParseAbsPath,
		readFile:  ioutil.ReadFile,
	}
}

// Visit scans the repo's OSM Ops target directory recursively, calling the
// specified visitor with the content of each OSM Git Ops file found.
// For now the only kind of Git Ops file OSM Ops can process is a file
// containing KduNsAction YAML. Any I/O errors that happen while traversing
// the target directory tree get collected in the returned error buffer as
// VisitErrors. Ditto for I/O errors that happen when reading or validating
// a Git Ops file as well as any error returned by the visitor.
func (k *KduNsActionRepoScanner) Visit(visitor KduNsActionProcessor) []error {
	es := []error{}
	filepath.Walk(k.targetDir.Value(), // (*)
		k.visitAllAndCollectErrors(visitor, &es))
	return es

	// (*) b/c targetDir is absolute, so is the path parameter passed to
	// the lambda returned by visitAllAndCollectErrors---see Walk docs.
}

func (k *KduNsActionRepoScanner) visitAllAndCollectErrors(
	visitor KduNsActionProcessor, acc *[]error) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			*acc = appendVisitError(path, err, *acc)
			return nil
		}
		if !k.isGitOpsFile(info) {
			return nil
		}
		if err := k.visitFile(path, visitor); err != nil {
			*acc = appendVisitError(path, err, *acc)
		}
		return nil
	}
}

func appendVisitError(path string, err error, errors []error) []error {
	visitError := &VisitError{AbsPath: path, Err: err}
	return append(errors, visitError)
}

func (k *KduNsActionRepoScanner) isGitOpsFile(info fs.FileInfo) bool {
	if !info.IsDir() {
		for _, ext := range k.fileExt {
			name := strings.ToLower(info.Name())
			if strings.HasSuffix(name, ext.Value()) {
				return true
			}
		}
	}
	return false
}

func (k *KduNsActionRepoScanner) visitFile(path string,
	visitor KduNsActionProcessor) error {
	var err error
	file := &KduNsActionFile{}

	absPath, err := k.parsePath(path)
	if err != nil { // (*)
		return err
	}
	file.FilePath = absPath

	yaml, err := k.readFile(path)
	if err != nil {
		return err
	}
	content, err := readKduNsAction(yaml)
	if err != nil {
		return err
	}
	file.Content = content

	return visitor.process(file)

	// (*) paranoia, it should never happen, path is already absolute
	// since targetDir is---see Walk docs.
}
