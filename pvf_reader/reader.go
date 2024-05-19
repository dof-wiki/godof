package pvf_reader

import (
	"bytes"
	"fmt"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"log"
	"os"
	"strings"
	"time"
)

type Unit struct {
	Type  uint8
	Value any
}

func (u *Unit) GetStr() string {
	return u.Value.(string)
}

type PvfReader struct {
	path   string
	encode string

	header      *PvfHeader
	f           *os.File
	stringTable *StringTable
	nString     *Lst
}

func NewPvfReader(path string, encode string) *PvfReader {
	return &PvfReader{
		path:   path,
		encode: encode,
		header: newPvfHeader(),
	}
}

func (p *PvfReader) readHeader() error {
	return p.header.read(p.f)
}

func (p *PvfReader) loadStringTable() error {
	st := time.Now()
	defer func() {
		fmt.Printf("load string table cost %s\n", time.Since(st))
	}()
	b, err := p.readFileContent("stringtable.bin")
	if err != nil {
		return err
	}
	p.stringTable = newStringTable(b, p.encode)
	return nil
}

func (p *PvfReader) loadNString() error {
	st := time.Now()
	defer func() {
		fmt.Printf("load nstring cost %s\n", time.Since(st))
	}()
	b, err := p.readFileContent("n_string.lst")
	if err != nil {
		return err
	}
	p.nString = NewLst(b, p.stringTable)
	return nil
}

func (p *PvfReader) Read() error {
	f, err := os.Open(p.path)
	if err != nil {
		return err
	}
	p.f = f
	if err := p.readHeader(); err != nil {
		return err
	}
	if err := p.loadStringTable(); err != nil {
		return err
	}
	if err := p.loadNString(); err != nil {
		return err
	}
	return nil
}

func (p *PvfReader) Close() {
	if p.f != nil {
		_ = p.f.Close()
		p.f = nil
	}
}

func (p *PvfReader) readBytes(offset, length uint32) ([]byte, error) {
	realOffset := p.header.ContentStartIdx + int64(offset)
	result := make([]byte, length)
	if _, err := p.f.ReadAt(result, realOffset); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PvfReader) ReadFileContent(path string) ([]*Unit, error) {
	c, err := p.readFileContent(path)
	if err != nil {
		return nil, err
	}
	return p.parseFileContent(c, ""), nil
}

func (p *PvfReader) ReadCommonFile(path string) (*CommonParser, error) {
	units, err := p.ReadFileContent(path)
	if err != nil {
		return nil, err
	}
	return NewCommonParser(units), nil
}

func (p *PvfReader) ReadLstFile(path string) (*Lst, error) {
	c, err := p.readFileContent(path)
	if err != nil {
		return nil, err
	}
	return NewLst(c, p.stringTable), nil
}

func (p *PvfReader) readFileContent(path string) ([]byte, error) {
	path = strings.ReplaceAll(path, "\\", "/")
	path, _ = strings.CutPrefix(strings.ToLower(path), "/")
	if fn, ok := p.header.FileTree[path]; !ok {
		return nil, os.ErrExist
	} else {
		c, err := p.readBytes(fn.RelativeOffset, fn.Length)
		if err != nil {
			log.Printf("read bytes err %v", err)
			return nil, err
		}
		c = binary_helper.DecryptCrc(c, fn.Crc32)
		return c, nil
	}
}

func (p *PvfReader) parseFileContent(c []byte, stringQuote string) []*Unit {
	unitNum := (len(c) - 2) / 5
	buffer := bytes.NewBuffer(c)
	shiftData := make([]byte, 2)
	units := make([]*Unit, 0, unitNum)
	_, _ = buffer.Read(shiftData)
	values := make([]any, 0, unitNum)
	types := make([]uint8, 0, unitNum)
	for i := 0; i < unitNum; i++ {
		var unitType uint8
		binary_helper.ReadAny(buffer, &unitType)
		var value any
		if unitType == 4 {
			var f float32
			binary_helper.ReadAny(buffer, &f)
			value = f
		} else {
			var v uint32
			binary_helper.ReadAny(buffer, &v)
			value = v
		}
		types = append(types, unitType)
		values = append(values, value)
	}
	for i := 0; i < unitNum; i++ {
		unitType := types[i]
		value := values[i]
		switch unitType {
		case 5, 6, 8:
			value, _ = p.stringTable.Get(value.(uint32))
		case 7:
			s, _ := p.stringTable.Get(value.(uint32))
			value = stringQuote + s + stringQuote
		case 9:
			path := strings.ToLower(p.nString.Get(value.(uint32)))
			strC, err := p.readFileContent(path)
			if err != nil {
				log.Printf("read file content err %v", err)
				continue
			}
			parser, err := NewStrParser(strC, p.encode)
			if err != nil {
				log.Printf("new str parser err %v", err)
				continue
			}
			st, err := p.stringTable.Get(values[i+1].(uint32))
			if err != nil {
				log.Printf("get string table err %v", err)
				continue
			}
			value = parser.Get(st)
		}
		units = append(units, &Unit{
			Type:  unitType,
			Value: value,
		})
	}
	return units
}
