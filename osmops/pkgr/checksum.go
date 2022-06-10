package pkgr

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/fluxcd/source-watcher/osmops/util/bytez"
	"github.com/fluxcd/source-watcher/osmops/util/file"
)

const ChecksumFileName = "checksums.txt"

func md5string(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

func computeCheckSum(target file.AbsPath) (string, error) {
	content, err := os.ReadFile(target.Value())
	if err != nil {
		return "", err
	}
	return md5string(content), nil
}

func writeCheckSumFileContent(archiveBaseName string, m pathToHash) io.Reader {
	buf := bytez.NewBuffer()
	archiveRelPaths := sortArchiveRelPaths(m)
	for _, relPath := range archiveRelPaths {
		baseNamePlusPath := path.Join(archiveBaseName, relPath)
		hash := m[relPath]
		line := fmt.Sprintf("%s\t%s\n", hash, baseNamePlusPath)
		io.Copy(buf, strings.NewReader(line))
	}
	return buf
}

func sortArchiveRelPaths(m pathToHash) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
