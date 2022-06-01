package file

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TODO. The path tests will probably fail on Windows since we're using
// Unix paths. We could use filepath.Join to make most of them platform
// independent but I'm not sure how to make absolute paths though...

var invalidPathFixtures = []string{"", " ", "\n", "\t "}

func TestInvalidPath(t *testing.T) {
	for k, d := range invalidPathFixtures {
		if err := IsStringPath(d); err == nil {
			t.Errorf("[%d] want: invalid; got: valid", k)
		}
	}
}

var parsePathFixtures = []struct {
	in   string
	want string
	rel  bool
}{
	{"/a/b/s", "/a/b/s", false}, {"r/e/l", "/r/e/l", true},
}

func TestParsePath(t *testing.T) {
	for k, d := range parsePathFixtures {
		if p, err := ParseAbsPath(d.in); err != nil {
			t.Errorf("[%d] want: valid parse; got: %v", k, err)
		} else {
			if !d.rel && d.want != p.Value() {
				t.Errorf("[%d] want: %s; got: %s", k, d.want, p.Value())
			}
			if d.rel && !strings.HasSuffix(p.Value(), d.want) {
				t.Errorf("[%d] want suffix: %s; got: %s", k, d.want, p.Value())
			}
		}
	}
}

var joinPathFixtures = []struct {
	base string
	rel  string
	want string
}{
	{"/a", "", "/a"}, {"/a/", " ", "/a"}, {"/a", "\t", "/a"},
	{"/a/", "b ", "/a/b"}, {"/a", "b\n", "/a/b"}, {"/a/b", "//c", "/a/b/c"},
}

func TestJoinPath(t *testing.T) {
	for k, d := range joinPathFixtures {
		if base, err := ParseAbsPath(d.base); err != nil {
			t.Errorf("[%d] want: valid parse; got: %v", k, err)
		} else {
			joined := base.Join(d.rel)
			if joined.Value() != d.want {
				t.Errorf("[%d] want: %s; got: %s", k, d.want, joined)
			}
		}
	}
}

func TestIsDir(t *testing.T) {
	if pwd, err := ParseAbsPath("."); err != nil {
		t.Errorf("want: valid parse; got: %v", err)
	} else {
		if err := pwd.IsDir(); err != nil {
			t.Errorf("want: pwd is a directory; got: %v", err)
		}

		notThere := pwd.Join("notThere")
		if err := notThere.IsDir(); err == nil {
			t.Errorf("want: not a directory; got directory: %v", notThere)
		}

		if tempFile, err := ioutil.TempFile("", "prefix"); err != nil {
			t.Errorf("couldn't create temp file: %v", err)
		} else {
			defer os.Remove(tempFile.Name())

			if tf, err := ParseAbsPath(tempFile.Name()); err != nil {
				t.Errorf("want: valid temp file parse; got: %v", err)
			} else {
				if err := tf.IsDir(); err == nil {
					t.Errorf("want: not a dir; got dir: %v", tf)
				}
			}
		}
	}
}
