package util

import (
	"fmt"
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

func TestEmptyString(t *testing.T) {
	if _, err := NewNonEmptyStr(""); err == nil {
		t.Errorf("instantiated a non-empty string with an empty string!")
	}
}

var nonEmptyStringFixtures = []string{" ", "\n", " wada wada "}

func TestNonEmptyString(t *testing.T) {
	for k, d := range nonEmptyStringFixtures {
		if s, err := NewNonEmptyStr(d); err != nil {
			t.Errorf("[%d] want: valid; got: %v", k, err)
		} else {
			if d != s.Value() {
				t.Errorf("[%d] want: %s; got: %s", k, d, s.Value())
			}
		}
	}
}

var invalidHostnameFixtures = []string{
	"", "\n", ":", ":80", "some.host:", "some host", "some host.com",
	"what?is.this", "em@il", "what.the.h*ll",
	"x1234567890123456789012345678901234567890123456789012345678901234.com",
}

func TestInvalidHostname(t *testing.T) {
	for k, d := range invalidHostnameFixtures {
		if err := IsHostname(d); err == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var validHostnameFixtures = []string{
	"::123", "1.2.3.4", "_h.com", "a-b.some_where", "some.host",
	"x12345678901234567890123456789012345678901234567890123456789012.com",
}

func TestValidHostname(t *testing.T) {
	for k, d := range validHostnameFixtures {
		if err := IsHostname(d); err != nil {
			t.Errorf("[%d] want: valid; got: %v", k, err)
		}
	}
}

var invalidHostnameAndPortFixtures = []string{
	"", "\n", ":", ":80", "some.host:", "some host:80", "some.host:123456789",
}

func TestInvalidHostnameAndPort(t *testing.T) {
	for k, d := range invalidHostnameAndPortFixtures {
		if err := IsHostAndPort(d); err == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var parseHostAndPortFixtures = []struct {
	in       string
	wantHost string
	wantPort int
}{
	{"h:0", "h", 0}, {"h:1", "h", 1}, {"h:65535", "h", 65535},
	{"[::123]:0", "::123", 0}, {"[::123]:1", "::123", 1},
	{"[::123]:65535", "::123", 65535},
	{"1.2.3.4:0", "1.2.3.4", 0}, {"1.2.3.4:1", "1.2.3.4", 1},
	{"1.2.3.4:65535", "1.2.3.4", 65535},
}

func TestParseHostAndPort(t *testing.T) {
	for k, d := range parseHostAndPortFixtures {
		if hp, err := ParseHostAndPort(d.in); err != nil {
			t.Errorf("[%d] want: valid parse; got: %v", k, err)
		} else {
			if d.wantHost != hp.Host() || d.wantPort != hp.Port() {
				t.Errorf("[%d] want: %s:%d; got: %v",
					k, d.wantHost, d.wantPort, hp)
			}

			repr := fmt.Sprintf("%s:%d", d.wantHost, d.wantPort)
			if repr != hp.String() {
				t.Errorf("[%d] want string repr: %s; got: %v", k, repr, hp)
			}
		}
	}
}

var httpUrlErrorFixtures = []string{"", "a", "a/b"}

func TestHttpUrlError(t *testing.T) {
	hp, _ := ParseHostAndPort("x:80")
	for k, d := range httpUrlErrorFixtures {
		if got, err := hp.Http(d); err == nil {
			t.Errorf("[%d] want error; got: %v", k, got)
		}
	}
}

var httpUrlFixtures = []struct {
	inPath string
	want   string
}{
	{"/", "http://x:80/"},
	{"/a", "http://x:80/a"}, {"/a/", "http://x:80/a/"},
	{"/a/b", "http://x:80/a/b"}, {"/a/b/", "http://x:80/a/b/"},
}

func TestHttpUrl(t *testing.T) {
	hp, _ := ParseHostAndPort("x:80")
	for k, d := range httpUrlFixtures {
		got, err := hp.Http(d.inPath)
		if err != nil {
			t.Fatalf("[%d] want string repr: %s; got: %v", k, d.want, err)
		}
		if got.String() != d.want {
			t.Errorf("[%d] want string repr: %s; got: %v", k, d.want, got)
		}
	}
}

var httpsUrlFixtures = []struct {
	inPath string
	want   string
}{
	{"/", "https://x:80/"},
	{"/a", "https://x:80/a"}, {"/a/", "https://x:80/a/"},
	{"/a/b", "https://x:80/a/b"}, {"/a/b/", "https://x:80/a/b/"},
}

func TestHttspUrl(t *testing.T) {
	hp, _ := ParseHostAndPort("x:80")
	for k, d := range httpsUrlFixtures {
		got, err := hp.Https(d.inPath)
		if err != nil {
			t.Fatalf("[%d] want string repr: %s; got: %v", k, d.want, err)
		}
		if got.String() != d.want {
			t.Errorf("[%d] want string repr: %s; got: %v", k, d.want, got)
		}
	}
}

func TestEmptyStrEnum(t *testing.T) {
	e := NewStrEnum()
	if e.IndexOf("") != NotALabel || e.IndexOf("x") != NotALabel {
		t.Errorf("empty enum should have no label indexes")
	}
	if e.LabelOf(0) != "" || e.LabelOf(1) != "" {
		t.Errorf("empty enum should have no labels")
	}
	if e.Validate("") == nil || e.Validate("x") == nil {
		t.Errorf("empty enum should always fail validation")
	}
}

type enumTest = struct {
	StrEnum
	A, B, C EnumIx
}

func NewEnumTest() enumTest {
	return enumTest{
		StrEnum: NewStrEnum("A", "b", "C"),
		A:       0,
		B:       1,
		C:       2,
	}
}

func TestStrEnumLookup(t *testing.T) {
	e := NewEnumTest()
	ixs := []EnumIx{e.A, e.B, e.C}
	for _, ix := range ixs {
		lbl := e.LabelOf(ix)
		if ix != e.IndexOf(lbl) {
			t.Errorf("want: %d == IndexOf(LabelOf(%d)); "+
				"got: %d != IndexOf(%s = LabelOf(%d)) == %d",
				ix, ix, ix, lbl, ix, e.IndexOf(lbl))
		}
	}
}

func TestStrEnumValidation(t *testing.T) {
	e := NewEnumTest()
	if err := e.Validate(e.LabelOf(e.A)); err != nil {
		t.Errorf("[1] want: valid; got: %v", err)
	}
	if err := e.Validate("wada wada"); err == nil {
		t.Errorf("[2] want: error; got: valid")
	}
}

func TestStrEnumCaseInsensitive(t *testing.T) {
	e := NewEnumTest()
	if err := e.Validate("B"); err != nil {
		t.Errorf("want: uppercase B is valid; got: %v", err)
	}
	if e.IndexOf("B") == NotALabel {
		t.Errorf("want: uppercase B is index of b; got: not a label")
	}
}
