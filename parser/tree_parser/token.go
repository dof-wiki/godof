package tree_parser

import (
	"fmt"
	"github.com/samber/lo"
	"strconv"
)

type TokenType int

var DelimiterChars = []string{"\t", "\n", " ", "\r"}

const (
	TokenNone   TokenType = iota
	TokenIgnore           // 注释
	TokenKey
	TokenString
	TokenInt
	TokenFloat
	TokenCommand
	TokenDelimiter
	TokenRawContent
)

type Token struct {
	tp      TokenType
	content string
}

func (t *Token) Render() string {
	switch t.tp {
	case TokenIgnore:
		return fmt.Sprintf("#%s\n", t.content)
	case TokenKey:
		return fmt.Sprintf("[%s]", t.content)
	case TokenString:
		return fmt.Sprintf("`%s`", t.content)
	case TokenCommand:
		return fmt.Sprintf("{%s}", t.content)
	default:
		return t.content
	}
}

func (t *Token) IsKey() bool {
	return t.tp == TokenKey
}

func (t *Token) IsCloseKey() bool {
	return t.tp == TokenKey && t.content[0] == '/'
}

func (t *Token) IsCloseKeyBy(key string) bool {
	return t.IsCloseKey() && t.content[1:] == key
}

func (t *Token) IsDelimiter() bool {
	return t.tp == TokenDelimiter
}

func (t *Token) IsIgnore() bool {
	return t.tp == TokenIgnore
}

func (t *Token) IsString() bool {
	return t.tp == TokenString
}

func (t *Token) IsInt() bool {
	return t.tp == TokenInt
}

func (t *Token) IsFloat() bool {
	return t.tp == TokenFloat
}

func (t *Token) RawContent() string {
	return t.content
}

func (t *Token) Clear() {
	t.content = ""
	t.tp = TokenNone
}

func (t *Token) Copy() *Token {
	return &Token{
		tp:      t.tp,
		content: t.content,
	}
}

func (t *Token) GetInt() int {
	return lo.Must(strconv.Atoi(t.content))
}

func (t *Token) GetFloat() float64 {
	return lo.Must(strconv.ParseFloat(t.content, 64))
}

func (t *Token) Value() any {
	switch t.tp {
	case TokenString:
		return t.content
	case TokenInt:
		return t.GetInt()
	case TokenFloat:
		return t.GetFloat()
	default:
		return t.content
	}
}

func NewDelimiterToken(c string) *Token {
	return &Token{
		tp:      TokenDelimiter,
		content: c,
	}
}

func NewRawContentToken(c string) *Token {
	return &Token{
		tp:      TokenRawContent,
		content: c,
	}
}

func GenTokens(val any) []*Token {
	tokens := make([]*Token, 0)
	switch v := val.(type) {
	case string:
		if lo.Contains(DelimiterChars, v) {
			tokens = append(tokens, &Token{tp: TokenDelimiter, content: v})
		} else {
			tokens = append(tokens, &Token{tp: TokenString, content: v})
		}
	case float64:
		tokens = append(tokens, &Token{tp: TokenFloat, content: strconv.FormatFloat(v, 'f', -1, 64)})
	case int:
		tokens = append(tokens, &Token{tp: TokenInt, content: strconv.Itoa(v)})
	case []string:
		for _, t := range v {
			tokens = append(tokens, &Token{tp: TokenString, content: t}, NewDelimiterToken("\t"))
		}
	case []int:
		for _, t := range v {
			tokens = append(tokens, &Token{tp: TokenInt, content: strconv.Itoa(t)}, NewDelimiterToken("\t"))
		}
	}
	return tokens
}

func GenTokenList(val ...any) []*Token {
	result := make([]*Token, 0)
	for _, v := range val {
		result = append(result, GenTokens(v)...)
		result = append(result, NewDelimiterToken("\t"))
	}
	return result
}
