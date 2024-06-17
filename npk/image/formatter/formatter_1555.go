package formatter

import (
	"bytes"
	"encoding/binary"
)

type Formatter1555 struct {
}

func (f *Formatter1555) ToRaw(data []byte) []byte {
	reader := bytes.NewReader(data)
	buf := make([]byte, len(data))
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

func (f *Formatter1555) formatColor(v1, v2 uint8) []uint8 {
	b := v1 & 0x1f
	b = b<<3 | b>>2
	g := (v1 >> 5) | (v2&3)<<3
	g = g<<3 | g>>2
	r := (v2 >> 2) & 0x1f
	r = r<<3 | r>>2
	a := v2 >> 7
	a = a * 0xff
	return []uint8{r, g, b, a}
}
