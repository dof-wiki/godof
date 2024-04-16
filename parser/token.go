package parser

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
	TokenNumber
	TokenDelimiter
)

type Token struct {
	tp      TokenType
	content string
}

func (t *Token) String() string {
	return fmt.Sprintf("(%d)%s", t.tp, t.content)
}

// Render 渲染成文本
func (t *Token) Render() string {
	switch t.tp {
	case TokenIgnore:
		return fmt.Sprintf("#%s\n", t.content)
	case TokenKey:
		return fmt.Sprintf("[%s]", t.content)
	case TokenString:
		return fmt.Sprintf("`%s`", t.content)
	default:
		return t.content
	}
}

func (t *Token) IsCloseKey() bool {
	return t.tp == TokenKey && t.content[0] == '/'
}

func (t *Token) IsCloseKeyBy(key string) bool {
	return t.tp == TokenKey && t.content[0] == '/' && t.content[1:] == key
}

func (t *Token) IsKey() bool {
	return t.tp == TokenKey && t.content[0] != '/'
}

func (t *Token) Type() TokenType {
	return t.tp
}

func (t *Token) IsString() bool {
	return t.tp == TokenString
}

func (t *Token) IsNumber() bool {
	return t.tp == TokenNumber
}

func (t *Token) IsDelimiter() bool {
	return t.tp == TokenDelimiter
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

func (t *Token) RawContent() string {
	return t.content
}

type TokenValue struct {
	cleanTokens     []*Token
	frontDelimiters []*Token
	backDelimiters  []*Token
}

func NewTokenValue(tokens []*Token) *TokenValue {
	tv := &TokenValue{
		cleanTokens:     make([]*Token, 0),
		frontDelimiters: make([]*Token, 0),
		backDelimiters:  make([]*Token, 0),
	}
	var start, end int
	for i, t := range tokens {
		if t.tp == TokenDelimiter {
			tv.frontDelimiters = append(tv.frontDelimiters, t)
		} else {
			start = i
			break
		}
	}
	for i := len(tokens) - 1; i >= 0; i-- {
		t := tokens[i]
		if t.tp == TokenDelimiter {
			tv.backDelimiters = append(tv.backDelimiters, t)
		} else {
			end = i
			break
		}
	}
	if start > end {
		tv.backDelimiters = make([]*Token, 0)
		return tv
	}
	tv.backDelimiters = lo.Reverse(tv.backDelimiters)
	tv.cleanTokens = tokens[start : end+1]
	return tv
}

func (tv *TokenValue) GetCleanTokens() []*Token {
	return tv.cleanTokens
}

func (tv *TokenValue) ReplaceValue(tokens []*Token) {
	tv.cleanTokens = tokens
}

func (tv *TokenValue) GetFull() []*Token {
	return append(tv.frontDelimiters, append(tv.cleanTokens, tv.backDelimiters...)...)
}

func (tv *TokenValue) GetStrings() ([]string, error) {
	return lo.FilterMap(tv.cleanTokens, func(item *Token, _ int) (string, bool) {
		if item.tp == TokenDelimiter {
			return "", false
		}
		return item.RawContent(), true
	}), nil
}

func (tv *TokenValue) GetString() (string, error) {
	list, err := tv.GetStrings()
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", &ErrEmptyValue{}
	}
	if len(list) > 1 {
		return "", &ErrValueType{}
	}
	return list[0], nil
}

func (tv *TokenValue) GetInts() ([]int, error) {
	var err error
	ret := lo.FilterMap(tv.cleanTokens, func(item *Token, _ int) (int, bool) {
		if item.tp == TokenDelimiter {
			return 0, false
		}
		v, er := strconv.Atoi(item.RawContent())
		if er != nil {
			err = &ErrValueType{}
		}
		return v, true
	})
	return ret, err
}

func (tv *TokenValue) GetInt() (int, error) {
	list, err := tv.GetInts()
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, &ErrEmptyValue{}
	}
	if len(list) > 1 {
		return 0, &ErrValueType{}
	}
	return list[0], nil
}

func NewDelimiterToken(c string) *Token {
	return &Token{
		tp:      TokenDelimiter,
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
	case int:
		tokens = append(tokens, &Token{tp: TokenNumber, content: strconv.Itoa(v)})
	case []string:
		for _, t := range v {
			tokens = append(tokens, &Token{tp: TokenString, content: t}, NewDelimiterToken("\t"))
		}
	case []int:
		for _, t := range v {
			tokens = append(tokens, &Token{tp: TokenNumber, content: strconv.Itoa(t)}, NewDelimiterToken("\t"))
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
