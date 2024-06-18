package formatter

import (
	"bytes"
	"encoding/binary"
	"io"
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

func (f *Formatter8888) ToRawCrop(data []byte, w, left, top, right, bottom int32) []byte {
	reader := bytes.NewReader(data)
	buf := make([]byte, 0, len(data))
	writer := bytes.NewBuffer(buf)
	for y := top; y < bottom; y++ {
		o := y * w * 2
		for x := left; x < right; x++ {
			reader.Seek(int64(o+x*2), io.SeekStart)
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
	}
	return writer.Bytes()
}
