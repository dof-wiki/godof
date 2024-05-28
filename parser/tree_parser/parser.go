package tree_parser

import (
	"strings"
)

type TreeParser struct {
	rawContent string
	closedKeys map[string]bool
	root       *TreeNode
}

func NewTreeParser(content string) *TreeParser {
	p := &TreeParser{
		rawContent: content,
		closedKeys: make(map[string]bool),
	}
	p.parse()
	return p
}

func (p *TreeParser) parse() {
	tokens := make([]*Token, 0)
	curToken := new(Token)
	for _, c := range p.rawContent {
		switch curToken.tp {
		case TokenNone:
			{
				switch c {
				case '#':
					curToken.tp = TokenIgnore
				case '[':
					curToken.tp = TokenKey
				case '`':
					curToken.tp = TokenString
				case '\t', ' ', '\n', '\r':
					tokens = append(tokens, &Token{
						tp:      TokenDelimiter,
						content: string(c),
					})
					continue
				case '{':
					curToken.tp = TokenCommand
				default:
					if (c >= '0' && c <= '9') || c == '-' {
						curToken.tp = TokenInt
						curToken.content += string(c)
					}
				}
			}
		case TokenIgnore:
			if c == '\n' {
				tokens = append(tokens, curToken.Copy())
				curToken.Clear()
			} else {
				curToken.content += string(c)
			}
		case TokenKey:
			if c == ']' {
				tokens = append(tokens, curToken.Copy())
				if curToken.RawContent()[0] == '/' {
					p.closedKeys[curToken.RawContent()[1:]] = true
				}
				curToken.Clear()
			} else {
				curToken.content += string(c)
			}
		case TokenString:
			if c == '`' {
				tokens = append(tokens, curToken.Copy())
				curToken.Clear()
			} else {
				curToken.content += string(c)
			}
		case TokenCommand:
			if c == '}' {
				tokens = append(tokens, curToken.Copy())
				curToken.Clear()
			} else {
				curToken.content += string(c)
			}
		case TokenInt, TokenFloat:
			if c == '.' {
				curToken.tp = TokenFloat
				curToken.content += string(c)
			} else if c < '0' || c > '9' {
				tokens = append(tokens, curToken.Copy(), NewDelimiterToken(string(c)))
				curToken.Clear()
			} else {
				curToken.content += string(c)
			}
		}
	}
	if curToken.tp != TokenNone {
		tokens = append(tokens, curToken.Copy())
		curToken.Clear()
	}
	p.root = &TreeNode{
		Label:    "root",
		IsClose:  false,
		Value:    tokens,
		children: nil,
	}
	p.root.parseChildren(p.closedKeys)
}

func (p *TreeParser) GetRoot() *TreeNode {
	return p.root
}

func (p *TreeParser) Render() string {
	result := "#PVF_File\n\n" + p.root.Render()
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result
}
