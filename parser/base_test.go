package parser

import (
	"fmt"
	"github.com/samber/lo"
	"os"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	content, err := os.ReadFile("../test/102010612.equ")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(string(content))
	s, err := p.GetString("name")
	if err != nil {
		t.Fatal(err)
	}
	if s != "苍穹幕落臂铠" {
		t.Fatal("name error")
	}
	s1, err := p.GetString("basic explain")
	t.Log(s1)
	p.SetString("basic explain", "aaaa")
	s2, err := p.GetString("basic explain")
	if err != nil {
		t.Fatal(err)
	}

	if s2 != "aaaa" {
		t.Fatal("basic explain set error")
	}

	s3, err := p.GetStrings("usable job")
	if err != nil {
		t.Fatal(err)
	}
	if len(s3) != 2 || s3[0] != "[fighter]" || s3[1] != "[at fighter]" {
		t.Fatalf("usable job err %v", s3)
	}
	p.SetStrings("usable job", []string{"[aaa]", "[bbb]"})
	s4, err := p.GetStrings("usable job")
	if err != nil {
		t.Fatal(err)
	}
	if len(s4) != 2 || s4[0] != "[aaa]" || s4[1] != "[bbb]" {
		t.Fatal("set usable job err")
	}

	n1, err := p.GetInt("value")
	if err != nil {
		t.Fatal(err)
	}
	if n1 != 161600 {
		t.Fatal("value error")
	}

	p.SetInt("value", 123)
	n2, err := p.GetInt("value")
	if err != nil {
		t.Fatal(err)
	}
	if n2 != 123 {
		t.Fatal("set value error")
	}

	n3, err := p.GetInts("variation")
	if err != nil {
		t.Fatal(err)
	}
	if len(n3) != 2 || n3[0] != 51 || n3[1] != 0 {
		t.Fatal("variation error")
	}

	p.SetInts("variation", []int{1, 2})
	n4, err := p.GetInts("variation")
	if err != nil {
		t.Fatal(err)
	}
	if len(n4) != 2 || n4[0] != 1 || n4[1] != 2 {
		t.Fatal("variation set error")
	}

	tv := p.GetAny("layer variation")
	if len(lo.Filter(tv.GetCleanTokens(), func(item *Token, _ int) bool {
		return !item.IsDelimiter()
	})) != 2 {
		t.Fatal("layer variation err")
	}
	for _, v := range tv.GetCleanTokens() {
		if v.IsInt() && v.RawContent() != "2790" {
			t.Fatal("layer variation err")
		}
		if v.IsString() && v.RawContent() != "gauntletc" {
			t.Fatal("layer variation err")
		}
	}

	p.SetAny("layer variation", GenTokenList(1, "test"))
	tv = p.GetAny("layer variation")
	if len(lo.Filter(tv.GetCleanTokens(), func(item *Token, _ int) bool {
		return !item.IsDelimiter()
	})) != 2 {
		t.Fatal("layer variation err")
	}
	for _, v := range tv.GetCleanTokens() {
		if v.IsInt() && v.RawContent() != "1" {
			t.Fatal("layer variation err")
		}
		if v.IsString() && v.RawContent() != "test" {
			t.Fatal("layer variation err")
		}
	}

	p.SetInt("creation rate", 10)
	p.SetInts("static", []int{10, 11, 12}, true)
	t.Log(p.Render())
}

func TestParseSpeed(t *testing.T) {
	content, err := os.ReadFile("../test/102010612.equ")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	st := time.Now()
	for i := 0; i < 100; i++ {
		NewParser(c)
	}
	d := time.Since(st)
	t.Logf("spent %d ms", d.Milliseconds())
}

func TestParseSubField(t *testing.T) {
	content, err := os.ReadFile("../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	p := NewParser(c)
	tv := p.GetAny("dungeon")
	subTv := tv.GetSub("static data")
	ints, _ := subTv.GetInts()
	if len(ints) != 3 {
		t.Fatal("dungeon static data err")
	}
}

func TestParseFloat(t *testing.T) {
	content, err := os.ReadFile("../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	p := NewParser(c)
	f, err := p.GetFloat("warroom")
	if err != nil {
		t.Fatal(err)
	}
	if f != 1.2 {
		t.Fatal("warroom err")
	}

	p.SetFloat("warroom", 1.3)
	f, err = p.GetFloat("warroom")
	if err != nil {
		t.Fatal(err)
	}
	if f != 1.3 {
		t.Fatal("warroom err")
	}

	for _, item := range p.GetAny("level property").cleanTokens {
		t.Log(item)
	}
}

func TestParseRawContent(t *testing.T) {
	content, err := os.ReadFile("../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	p := NewParser(c)
	p.SetAny("test any", []*Token{
		NewRawContentToken("`test`\t989"),
	}, true)
	p.SetAny("test any 2", []*Token{
		NewRawContentToken("[sub key]\n`test`\t989\n[/sub key]"),
	}, true)
	fmt.Println(p.Render())
}

func Test2(t *testing.T) {
	content, err := os.ReadFile("../test/earthbreak.skl")
	if err != nil {
		t.Fatal(err)
	}
	c := string(content)
	p := NewParser(c)
	fmt.Println(p.Render())
}
