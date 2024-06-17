package formatter

import (
	"errors"
)

const (
	IMAGE_FORMAT_1555 = 14
	IMAGE_FORMAT_4444 = 15
	IMAGE_FORMAT_8888 = 16
)

type Formatter interface {
	ToRaw(data []byte) []byte
}

func FormatToRaw(data []byte, format int32) ([]byte, error) {
	var formatter Formatter
	switch format {
	case IMAGE_FORMAT_1555:
		formatter = new(Formatter1555)
	case IMAGE_FORMAT_4444:
		formatter = new(Formatter4444)
	case IMAGE_FORMAT_8888:
		formatter = new(Formatter8888)
	default:
		return nil, errors.New("unknown formatter")
	}
	return formatter.ToRaw(data), nil
}
