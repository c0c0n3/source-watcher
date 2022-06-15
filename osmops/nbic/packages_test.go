package nbic

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fluxcd/source-watcher/osmops/util/file"
)

func findTestDataDir(pkgDirName string) file.AbsPath {
	_, thisFileName, _, _ := runtime.Caller(1)
	enclosingDir := filepath.Dir(thisFileName)
	testDataDir := filepath.Join(enclosingDir, "packages_test_dir", pkgDirName)
	p, _ := file.ParseAbsPath(testDataDir)

	return p
}

func md5string(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

func callCreateOrUpdatePackage(times int, pkgDirName string) (*mockNbi, error) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)
	pkgSrc := findTestDataDir(pkgDirName)

	var err error
	for k := 0; k < times; k++ {
		err = nbic.CreateOrUpdatePackage(pkgSrc)
		if err != nil {
			break
		}
	}
	return nbi, err
}

func checkUploadedPackage(t *testing.T, mockNbi *mockNbi, req *http.Request,
	pkgDirName string) {
	gotFilename := req.Header.Get("Content-Filename")
	gotHash := req.Header.Get("Content-File-MD5")
	gotPkgTgzData := mockNbi.packages[pkgDirName]

	wantFilename := fmt.Sprintf("%s.tar.gz", pkgDirName)
	if gotFilename != wantFilename {
		t.Errorf("want file: %s; got: %s", wantFilename, gotFilename)
	}
	wantHash := md5string(gotPkgTgzData)
	if gotHash != wantHash {
		t.Errorf("want hash: %s; got: %s", wantHash, gotHash)
	}
}

func checkUnsupportedPackageErr(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("want err; got: nil")
	}
	if !strings.HasPrefix(err.Error(), "unsupported package type") {
		t.Errorf("want unsupported pkg err; got: %v", err)
	}
}

func runCreatePackageTest(t *testing.T, pkgDirName string) {
	mockNbi, err := callCreateOrUpdatePackage(1, pkgDirName)

	if err != nil {
		t.Errorf("want: create package; got: %v", err)
	}
	if len(mockNbi.exchanges) != 2 { // #1 = get token
		t.Fatalf("want: one req to create package; got: %d",
			len(mockNbi.exchanges)-1)
	}

	rr := mockNbi.exchanges[1]
	checkUploadedPackage(t, mockNbi, rr.req, pkgDirName)
	if rr.res.StatusCode != http.StatusCreated {
		t.Errorf("want status: %d; got: %d",
			http.StatusCreated, rr.res.StatusCode)
	}
}

func runUpdatePackageTest(t *testing.T, pkgDirName string) {
	mockNbi, err := callCreateOrUpdatePackage(2, pkgDirName)

	if err != nil {
		t.Errorf("want: update package; got: %v", err)
	}
	if len(mockNbi.exchanges) != 4 { // #1 = get token
		msg := "want: one initial req to create the package, then one failed " +
			"attempt to create it again, followed by one to update it; got: %d"
		t.Fatalf(msg, len(mockNbi.exchanges)-1)
	}

	failedCreateExchange := mockNbi.exchanges[2]
	checkUploadedPackage(t, mockNbi, failedCreateExchange.req, pkgDirName)
	if failedCreateExchange.res.StatusCode != http.StatusConflict {
		t.Errorf("want create status: %d; got: %d",
			http.StatusConflict, failedCreateExchange.res.StatusCode)
	}

	updateExchange := mockNbi.exchanges[3]
	checkUploadedPackage(t, mockNbi, updateExchange.req, pkgDirName)
	if updateExchange.res.StatusCode != http.StatusOK {
		t.Errorf("want update status: %d; got: %d",
			http.StatusOK, updateExchange.res.StatusCode)
	}
}

func TestCreateKnfPackage(t *testing.T) {
	runCreatePackageTest(t, "openldap_knf")
}

func TestCreateNsPackage(t *testing.T) {
	runCreatePackageTest(t, "openldap_ns")
}

func TestUpdateKnfPackage(t *testing.T) {
	runUpdatePackageTest(t, "openldap_knf")
}

func TestUpdateNsPackage(t *testing.T) {
	runUpdatePackageTest(t, "openldap_ns")
}

func TestPackErrOnSourceDirAccess(t *testing.T) {
	mockNbi, err := callCreateOrUpdatePackage(1, "not-there_knf")

	if _, ok := err.(*file.VisitError); !ok {
		t.Errorf("want: visit error; got: %v", err)
	}
	if len(mockNbi.exchanges) > 0 {
		t.Errorf("want: no req to create or update package; got: %d",
			len(mockNbi.exchanges))
	}
}

func TestCreateUnsupportedPackage(t *testing.T) {
	mockNbi, err := callCreateOrUpdatePackage(1, "unsupported")
	if len(mockNbi.exchanges) > 0 {
		t.Errorf("want: no req to create or update package; got: %d",
			len(mockNbi.exchanges))
	}
	checkUnsupportedPackageErr(t, err)
}

func TestUpdateUnsupportedPackage(t *testing.T) {
	nbi := newMockNbi()
	urls := newConn()
	nbic, _ := New(urls, usrCreds, nbi.exchange)
	pkgSrc := findTestDataDir("unsupported")
	reader, _ := newPkgReader(pkgSrc)
	handler := pkgHandler{
		session: nbic,
		pkg:     reader,
	}

	_, err := handler.update()
	checkUnsupportedPackageErr(t, err)
}
