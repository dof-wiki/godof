package binary_helper

import (
	"encoding/binary"
	"golang.org/x/exp/constraints"
	"io"
	"log"
)

func ReadInt32(f io.Reader) int32 {
	var length int32
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		log.Fatal(err)
	}
	return length
}

func ReadUInt32(f io.Reader) uint32 {
	var length uint32
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		log.Fatal(err)
	}
	return length
}

func ReadAny[T any](f io.Reader, i *T) {
	if err := binary.Read(f, binary.LittleEndian, i); err != nil {
		log.Fatal(err)
	}
}

func ReadStr(f io.Reader) (string, int32) {
	var length int32
	result := ReadBytesByLen(f, &length)
	return string(result), length
}

func ReadBytesByLen[T constraints.Integer](f io.Reader, length *T) []byte {
	err := binary.Read(f, binary.LittleEndian, length)
	if err != nil {
		if err == io.EOF {
			return make([]byte, 0)
		}
		log.Fatal(err)
	}
	buffer := make([]byte, *length)
	_, err = io.ReadFull(f, buffer)
	if err != nil {
		log.Fatal(err)
	}
	return buffer
}
