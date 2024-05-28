package tree_parser

import "github.com/samber/lo"

type Value []*Token

func (v Value) Get() []any {
	return lo.FilterMap(v, func(token *Token, _ int) (any, bool) {
		return token.Value(), !(token.IsDelimiter() || token.IsIgnore())
	})
}

func (v Value) GetInts() []int {
	return lo.FilterMap(v, func(token *Token, _ int) (int, bool) {
		if token.IsInt() {
			return token.GetInt(), true
		}
		return 0, false
	})
}

func (v Value) GetFloats() []float64 {
	return lo.FilterMap(v, func(token *Token, _ int) (float64, bool) {
		if token.IsFloat() {
			return token.GetFloat(), true
		}
		return 0, false
	})
}

func (v Value) GetStrings() []string {
	return lo.FilterMap(v, func(token *Token, _ int) (string, bool) {
		if token.IsString() {
			return token.content, true
		}
		return "", false
	})
}

func (v Value) GetString() string {
	list := v.GetStrings()
	if len(list) == 0 {
		return ""
	}
	return list[0]
}

func (v Value) GetInt() int {
	list := v.GetInts()
	if len(list) == 0 {
		return 0
	}
	return list[0]
}

func (v Value) GetFloat() float64 {
	list := v.GetFloats()
	if len(list) == 0 {
		return 0
	}
	return list[0]
}
