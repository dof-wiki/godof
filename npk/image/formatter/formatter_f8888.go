package formatter

import (
	"bytes"
	"encoding/binary"
)

type Formatter8888 struct {
}

func (f *Formatter8888) ToRaw(data []byte) []byte {
	reader := bytes.NewReader(data)
	buf := make([]byte, 0, len(data))
	writer := bytes.NewBuffer(buf)

	for {
		var b, g, r, a uint8
		err := binary.Read(reader, binary.LittleEndian, &b)
		if err != nil {
			break
		}
		err = binary.Read(reader, binary.LittleEndian, &g)
		if err != nil {
			break
		}
		err = binary.Read(reader, binary.LittleEndian, &r)
		if err != nil {
			break
		}
		err = binary.Read(reader, binary.LittleEndian, &a)
		if err != nil {
			break
		}
		writer.WriteByte(r)
		writer.WriteByte(g)
		writer.WriteByte(b)
		writer.WriteByte(a)
	}
	return writer.Bytes()
}
