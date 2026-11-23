package palette

import (
	"bytes"
	"testing"
)

func TestToCSSDeterministic(t *testing.T) {
	p, err := Generate("#1e1e2e")
	if err != nil {
		t.Fatal(err)
	}
	var a, b bytes.Buffer
	p.ToCSS(&a)
	p.ToCSS(&b)
	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatal("ToCSS output differs between renders")
	}
}
