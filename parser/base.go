package parser

type Parser struct {
	rawContent string
	tokens     []*Token
	closedKeys map[string]bool
}

func NewParser(content string) *Parser {
	p := &Parser{
		rawContent: content,
		tokens:     make([]*Token, 0),
		closedKeys: make(map[string]bool),
	}
	p.parse()
	return p
}

func (p *Parser) parse() {
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
	p.tokens = tokens
}

func (p *Parser) isClosedKey(key string) bool {
	return p.closedKeys[key]
}

func (p *Parser) Render() string {
	result := ""
	for _, t := range p.tokens {
		result += t.Render()
	}
	return result
}

func (p *Parser) GetTokens() []*Token {
	return p.tokens
}

func (p *Parser) getKV(key string) (*TokenValue, int, int) {
	for i, t := range p.tokens {
		if t.IsKey() && t.RawContent() == key {
			for j, t2 := range p.tokens[i+1:] {
				if p.isClosedKey(key) {
					if t2.IsCloseKeyBy(key) {
						return NewTokenValue(p.tokens[i+1 : i+1+j]), i + 1, i + j
					}
				} else if t2.IsKey() {
					return NewTokenValue(p.tokens[i+1 : i+1+j]), i + 1, i + j
				}
			}
			return NewTokenValue(p.tokens[i+1:]), i + 1, len(p.tokens)
		}
	}
	return nil, 0, 0
}

func (p *Parser) addKV(key string, values []*Token, isClosed ...bool) {
	newTokens := []*Token{
		NewDelimiterToken("\n"),
		{
			tp:      TokenKey,
			content: key,
		},
		NewDelimiterToken("\n"),
		NewDelimiterToken("\t"),
	}
	newTokens = append(newTokens, values...)
	newTokens = append(newTokens, NewDelimiterToken("\n"))
	if len(isClosed) > 0 && isClosed[0] {
		newTokens = append(newTokens, &Token{
			tp:      TokenKey,
			content: "/" + key,
		}, NewDelimiterToken("\n"))
		p.closedKeys[key] = true
	}
	p.tokens = append(p.tokens, newTokens...)
}

func (p *Parser) setKV(key string, values []*Token, isClosed ...bool) {
	tokens, i, j := p.getKV(key)
	if tokens == nil {
		p.addKV(key, values, isClosed...)
		return
	}
	tokens.ReplaceValue(values)
	// 把 p.tokens[i:j+1] 替换为 values
	p.tokens = append(p.tokens[:i], append(tokens.GetFull(), p.tokens[j+1:]...)...)
}

func (p *Parser) delKV(key string) {
	tv, i, j := p.getKV(key)
	if tv == nil {
		return
	}
	if p.isClosedKey(key) {
		j++
	}
	// 把 p.tokens[i-1:j+1] 删除
	p.tokens = append(p.tokens[:i-1], p.tokens[j+1:]...)
}

func (p *Parser) GetString(key string) (string, error) {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return "", &ErrEmptyValue{}
	}
	return tv.GetString()
}

func (p *Parser) SetString(key, value string, isClosed ...bool) {
	p.setKV(key, GenTokens(value), isClosed...)
}

func (p *Parser) GetInt(key string) (int, error) {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return 0, &ErrEmptyValue{}
	}
	return tv.GetInt()
}

func (p *Parser) SetInt(key string, value int, isClosed ...bool) {
	p.setKV(key, GenTokens(value), isClosed...)
}

func (p *Parser) GetInts(key string) ([]int, error) {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return nil, &ErrEmptyValue{}
	}
	return tv.GetInts()
}

func (p *Parser) SetInts(key string, value []int, isClosed ...bool) {
	p.setKV(key, GenTokens(value), isClosed...)
}

func (p *Parser) GetStrings(key string) ([]string, error) {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return nil, &ErrEmptyValue{}
	}
	return tv.GetStrings()
}

func (p *Parser) SetStrings(key string, value []string, isClosed ...bool) {
	p.setKV(key, GenTokens(value), isClosed...)
}

func (p *Parser) SetFloat(key string, value float64, isClosed ...bool) {
	p.setKV(key, GenTokens(value), isClosed...)
}

func (p *Parser) GetFloat(key string) (float64, error) {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return 0, &ErrEmptyValue{}
	}
	return tv.GetFloat()
}

func (p *Parser) GetAny(key string) *TokenValue {
	tv, _, _ := p.getKV(key)
	if tv == nil {
		return nil
	}
	return tv
}

func (p *Parser) SetAny(key string, value []*Token, isClosed ...bool) {
	p.setKV(key, value, isClosed...)
}

func (p *Parser) DelKey(key string) {
	p.delKV(key)
}

func (p *Parser) AddAny(key string, value []*Token, isClosed ...bool) {
	p.addKV(key, value, isClosed...)
}
