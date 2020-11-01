package filereader

import (
	"io"
	"testing"
)

func TestLineReader(t *testing.T) {

	r, err := NewLineReader("/Users/adamtom3/Documents/real_documents/test/esTool/pkg/streambuf/streambuf.go", 4096)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for {
		line, err := r.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Log(err)
			t.FailNow()
		}
		t.Log(line)
	}

}
