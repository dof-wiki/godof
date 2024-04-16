package parser

import (
	"os"
	"testing"
)

func TestLstParser(t *testing.T) {
	content, err := os.ReadFile("../test/test.lst")
	if err != nil {
		t.Fatal(err)
	}
	p := NewLstParser(string(content))

	if len(p.GetPathMap()) != 4 {
		t.Fatal("path map err")
	}

	if err := p.AddPath(888, "test/path/xx.equ"); err != nil {
		t.Fatal(err)
	}
	t.Log(p.Render())
}
