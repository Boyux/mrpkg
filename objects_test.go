package mrpkg

import (
	"bytes"
	"strconv"
	"testing"
)

func TestGetObjAndPutObj(t *testing.T) {
	for i := 0; i < 10; i++ {
		var buf bytes.Buffer
		buf.WriteString("hello world")
		PutObj(&buf)
	}

	buf := GetObj[*bytes.Buffer]()
	if buf.Len() != 0 {
		t.Logf("*bytes.Buffer: %s\n", strconv.Quote(buf.String()))
		if buf.String() != "hello world" {
			t.Errorf("error string from *bytes.Buffer: %s\n", strconv.Quote(buf.String()))
		}
	}
}
