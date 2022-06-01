package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AbsPath struct{ data string }

func (d AbsPath) Value() string {
	return d.data
}

func IsStringPath(value interface{}) error {
	s, _ := value.(string)
	_, err := ParseAbsPath(s)
	return err
}

func ParseAbsPath(path string) (AbsPath, error) {
	path = strings.TrimSpace(path) // (*)
	if len(path) == 0 {
		return AbsPath{},
			errors.New("must be a non-empty, non-whitespace-only string")
	}
	if p, err := filepath.Abs(path); err != nil {
		return AbsPath{}, err
	} else {
		return AbsPath{data: p}, nil
	}

	// (*) Abs doesn't trim space, e.g. Abs('/a/b ') == '/a/b '.
}

func (d AbsPath) Join(relativePath string) AbsPath {
	rest := strings.TrimSpace(relativePath) // (1)
	return AbsPath{
		data: filepath.Join(d.Value(), rest), // (2)
	}

	// (1) Join doesn't trim space, e.g. Join("/a", "/b ") == "/a/b "
	// (2) In principle this is wrong since we don't know if relativePath
	// is a valid path according to the FS we're running on. (Join doesn't
	// check that.) So we could potentially return an inconsistent AbsPath.
	// Go's standard lib is quite weak in the handling of abstract paths,
	// i.e. independent of OS, so this is the best we can do. See e.g.
	// - https://stackoverflow.com/questions/35231846
}

func (d AbsPath) IsDir() error {
	if f, err := os.Stat(d.Value()); err != nil {
		return err
	} else {
		if !f.IsDir() {
			return fmt.Errorf("not a directory: %v", d.Value())
		}
	}
	return nil
}
