package zlib_helper

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

func FillBytes(data []byte, size int) []byte {
	fillSize := size - len(data)
	if fillSize > 0 {
		zeroFill := make([]byte, fillSize)
		data = append(data, zeroFill...)
	}
	return data
}

/*
def zlib_decompress(data: bytes) -> bytes:
    if data.startswith(b'\x78'):
        data = zlib.decompress(data)
    else:
        header_index = data.rfind(b'\x78')
        try:
            data = zlib.decompress(data[header_index:] + data)
        except zlib.error:
            data = zlib.decompress(data[header_index:header_index + 2] + data)

    return data
*/

func decompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	if len(data) > 0 && data[0] == 0x78 {
		return decompress(data)
	}
	headerIndex := bytes.IndexByte(data, 0x78)
	if headerIndex == -1 {
		return nil, errors.New("no valid zlib header found")
	}
	d, err := decompress(append(data[headerIndex:], data...))
	if err != nil {
		d, err = decompress(append(data[headerIndex:headerIndex+2], data...))
	}
	return d, err
}
