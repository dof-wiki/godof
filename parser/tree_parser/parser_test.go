package tree_parser

import (
	"fmt"
	"os"
	"testing"
)

func TestTreeParser(t *testing.T) {
	content, err := os.ReadFile("../../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	p := NewTreeParser(string(content))
	fmt.Println(p.root.Render())
}

func TestParseRawContent(t *testing.T) {
	content, err := os.ReadFile("../../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	p := NewTreeParser(c)
	nameNode := p.GetRoot().GetFirstChild("name")
	name := nameNode.Value.GetString()
	fmt.Println(name)
	data := p.GetRoot().GetFirstChild("level property").Value.Get()
	fmt.Println(data)
	p.GetRoot().AddChild(NewTreeNode("test", false, GenTokenList("aaaaa", 222)...))
	p.GetRoot().AddChild(NewTreeNode("test any", true, NewRawContentToken("`test`\t989")))
	p.GetRoot().AddChild(NewTreeNode("test any 2", true, NewRawContentToken("[sub key]\n`test`\t989\n[/sub key]")))
	//fmt.Println(p.Render())
}
