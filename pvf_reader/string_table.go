package pvf_reader

import (
	"encoding/binary"
	"github.com/dof-wiki/godof/utils"
	"sync"
)

type StringTable struct {
	mu      sync.Mutex
	encode  string
	length  uint32
	content []byte
	chunk   []byte

	cached map[uint32]string
}

func newStringTable(buf []byte, encode string) *StringTable {
	s := &StringTable{
		encode: encode,
		cached: make(map[uint32]string),
	}
	s.length = binary.LittleEndian.Uint32(buf[:4])
	s.content = buf[4:]
	s.chunk = buf[4+s.length*4+4:]
	return s
}

func (s *StringTable) get(idx uint32) (string, error) {
	c := s.content[idx*4 : idx*4+8]
	idx1 := binary.LittleEndian.Uint32(c[0:4])
	idx2 := binary.LittleEndian.Uint32(c[4:8])
	bias := s.length*4 + 4
	ret := s.chunk[idx1-bias : idx2-bias]
	return utils.Decode(ret, s.encode)
}

func (s *StringTable) Get(idx uint32) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.cached[idx]; ok {
		return v, nil
	}
	v, err := s.get(idx)
	if err != nil {
		return "", err
	}
	s.cached[idx] = v
	return v, nil
}
