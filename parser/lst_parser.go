package parser

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"sort"
	"strconv"
	"strings"
)

type LstParser struct {
	*Parser
	pathMap map[int]string
}

func NewLstParser(content string) *LstParser {
	p := &LstParser{
		Parser:  NewParser(content),
		pathMap: map[int]string{},
	}
	p.parseLst()
	return p
}

func (p *LstParser) parseLst() {

	k := -1
	for _, t := range p.tokens {
		switch t.tp {
		case TokenNumber:
			k, _ = strconv.Atoi(t.content)
		case TokenString:
			if k >= 0 {
				p.pathMap[k] = t.content
				k = -1
			}
		default:
			continue
		}
	}
}

func (p *LstParser) GetPathMap() map[int]string {
	return p.pathMap
}

func (p *LstParser) GetPath(idx int) string {
	return p.pathMap[idx]
}

func (p *LstParser) AddPath(idx int, path string) error {
	if _, has := p.pathMap[idx]; has {
		return errors.New(fmt.Sprintf("idx %d repeated.", idx))
	}
	p.pathMap[idx] = path
	return nil
}

func (p *LstParser) DelPath(idx int) {
	delete(p.pathMap, idx)
}

func (p *LstParser) Render() string {
	result := []string{"#PVF_File", ""}
	keys := lo.Keys(p.pathMap)
	sort.Ints(keys)
	for _, idx := range keys {
		result = append(result, strconv.Itoa(idx))
		result = append(result, fmt.Sprintf("`%s`", p.pathMap[idx]))
	}
	return strings.Join(result, "\n")
}
