package utils

import (
	"bytes"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

func DecodeCP949(data []byte) (string, error) {
	// 使用 CP949 解码器
	decoder := korean.EUCKR.NewDecoder()

	// 创建一个 transformer 来处理字节到字符串的转换
	reader := transform.NewReader(strings.NewReader(string(data)), decoder)

	// 读取并转换数据
	decodedData, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedData), nil
}

func DecodeBig5(data []byte) (string, error) {
	// 使用 Big5 解码器
	decoder := traditionalchinese.Big5.NewDecoder()

	// 创建一个 transformer 来处理字节到字符串的转换
	reader := transform.NewReader(strings.NewReader(string(data)), decoder)

	// 读取并转换数据
	decodedData, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedData), nil
}

func Decode(data []byte, encode string) (string, error) {
	encode = strings.ReplaceAll(strings.ToLower(encode), "-", "")
	decoder := unicode.UTF8.NewDecoder()
	switch encode {
	case "big5":
		decoder = traditionalchinese.Big5.NewDecoder()
	case "cp949", "euc_kr":
		decoder = korean.EUCKR.NewDecoder()
	case "gbk":
		decoder = simplifiedchinese.GBK.NewDecoder()
	default:
		return string(data), nil
	}
	reader := transform.NewReader(strings.NewReader(string(data)), decoder)

	// 读取并转换数据
	decodedData, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedData), nil
}

func TrimStringZeros(s string) string {
	return string(TrimBytesZeros([]byte(s)))
}

func TrimBytesZeros(data []byte) []byte {
	idx := bytes.IndexByte(data, 0x00)
	if idx != -1 {
		return data[:idx]
	}
	return data
}
