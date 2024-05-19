package pvf_reader

import (
	"bytes"
	"github.com/dof-wiki/godof/utils"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"log"
	"strings"
	"sync"
)

type Lst struct {
	Data map[uint32]string
}

func NewLst(c []byte, stringTable *StringTable) *Lst {
	data := make(map[uint32]string)
	i := 2
	for {
		if i+10 >= len(c) {
			break
		}
		buf := bytes.NewBuffer(c[i : i+10])
		var a, b int8
		var aa, bb uint32
		binary_helper.ReadAny(buf, &a)
		binary_helper.ReadAny(buf, &aa)
		binary_helper.ReadAny(buf, &b)
		binary_helper.ReadAny(buf, &bb)
		var index, strIdx uint32
		if a == 2 {
			index = aa
		} else if a == 7 {
			strIdx = aa
		}
		if b == 2 {
			index = bb
		} else if b == 7 {
			strIdx = bb
		}

		s, err := stringTable.Get(strIdx)
		if err != nil {
			log.Println("err is ", err)
		}
		data[index] = s
		i += 10
	}

	lst := &Lst{
		Data: data,
	}
	return lst
}

func (lst *Lst) Get(idx uint32) string {
	return lst.Data[idx]
}

type StrParser struct {
	mu   sync.Mutex
	Data map[string]string
}

func NewStrParser(c []byte, encode string) (*StrParser, error) {
	s, err := utils.Decode(c, encode)
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		has := false
		for _, i := range line {
			if i == '>' {
				has = true
			}
		}
		if !has {
			continue
		}
		ret := strings.SplitN(line, ">", 1)
		if len(ret) == 2 {
			data[ret[0]] = ret[1]
		}
	}
	return &StrParser{
		Data: data,
	}, nil
}

func (p *StrParser) Get(idx string) string {
	return p.Data[idx]
}

type FieldValue []*Unit

func (fv FieldValue) GetInts() []uint32 {
	ret := make([]uint32, 0, len(fv))
	for _, u := range fv {
		if v, ok := u.Value.(uint32); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func (fv FieldValue) GetInt() uint32 {
	if len(fv) != 1 {
		return 0
	}
	return fv[0].Value.(uint32)
}

func (fv FieldValue) GetStrs() []string {
	ret := make([]string, 0, len(fv))
	for _, u := range fv {
		if v, ok := u.Value.(string); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func (fv FieldValue) GetStr() string {
	if len(fv) != 1 {
		return ""
	}
	return fv[0].GetStr()
}

func (fv FieldValue) GetFloats() []float32 {
	var ret []float32
	for _, u := range fv {
		if v, ok := u.Value.(uint32); ok {
			ret = append(ret, float32(v))
		}
		if v, ok := u.Value.(float32); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

type Field struct {
	Key         string
	Value       FieldValue
	SelfClosing bool
}

type CommonParser struct {
	Fields map[string][]*Field
	units  []*Unit
}

func NewCommonParser(unit []*Unit) *CommonParser {
	p := &CommonParser{
		Fields: make(map[string][]*Field),
		units:  unit,
	}
	p.parser()
	return p
}

func (p *CommonParser) parser() {
	closingKey := make(map[string]bool)
	for _, u := range p.units {
		if u.Type == UnitTypeKey && strings.HasPrefix(u.GetStr(), "[/") {
			s := u.GetStr()
			closingKey[s[2:len(s)-1]] = true
		}
	}
	var key string
	values := make([]*Unit, 0)
	for _, u := range p.units {
		switch u.Type {
		case UnitTypeKey:
			s := u.GetStr()
			if key == "" {
				key = s[1 : len(s)-1]
			} else {
				isClose := strings.HasPrefix(s, "[/")
				if closingKey[key] {
					if isClose && s[2:len(s)-1] == key {
						p.Fields[key] = append(p.Fields[key], &Field{
							Key:         key,
							Value:       values,
							SelfClosing: true,
						})
						values = make([]*Unit, 0)
						key = ""
					} else {
						values = append(values, u)
					}
				} else {
					p.Fields[key] = append(p.Fields[key], &Field{
						Key:   key,
						Value: values,
					})
					values = make([]*Unit, 0)
					if strings.HasPrefix(s, "[/") {
						key = ""
					} else {
						key = s[1 : len(s)-1]
					}
				}
			}
		default:
			values = append(values, u)
		}
	}
}

func (p *CommonParser) GetField(key string) []*Field {
	return p.Fields[key]
}
