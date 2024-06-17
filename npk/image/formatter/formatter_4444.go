package formatter

import (
	"bytes"
	"encoding/binary"
)

type Formatter4444 struct{}

func (f *Formatter4444) ToRaw(data []byte) []byte {
	reader := bytes.NewReader(data)
	buf := make([]byte, 0, len(data))
	writer := bytes.NewBuffer(buf)

	for {
		var byte1, byte2 uint8
		err := binary.Read(reader, binary.LittleEndian, &byte1)
		if err != nil {
			break
		}
		err = binary.Read(reader, binary.LittleEndian, &byte2)
		if err != nil {
			break
		}
		for _, v := range f.formatColor(byte1, byte2) {
			writer.WriteByte(v)
		}
	}
	return writer.Bytes()
}

func (f *Formatter4444) formatColor(v1, v2 uint8) []uint8 {
	b := (v1 & 0xf) << 4
	g := v1 & 0xf0
	r := (v2 & 0xf) << 4
	a := v2 & 0xf0
	return []uint8{r, g, b, a}
}
