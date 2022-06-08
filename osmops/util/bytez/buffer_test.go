package bytez

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

const dataSize = 4 * 1024

func makeData() []byte {
	data := make([]byte, dataSize)
	for k := 0; k < dataSize; k++ {
		data[k] = byte(k % 256)
	}
	return data
}

func checkData(t *testing.T, data []byte) {
	if len(data) != dataSize {
		t.Errorf("want size: %d; got: %d", dataSize, len(data))
	}
	for k := 0; k < len(data); k++ {
		want := byte(k % 256)
		if data[k] != want {
			t.Errorf("[%d] want: %d; got: %d", k, want, data[k])
		}
	}
}

func writeAll(dest io.WriteCloser) {
	defer dest.Close()
	src := bytes.NewBuffer(makeData())
	io.Copy(dest, src)
}

func readAll(src io.ReadCloser) []byte {
	defer src.Close()
	data, _ := ioutil.ReadAll(src)
	return data
}

func TestWriteThenRead(t *testing.T) {
	buf := NewBuffer()
	writeAll(buf)
	got := readAll(buf)
	checkData(t, got)
}
