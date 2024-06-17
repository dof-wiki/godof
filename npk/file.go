package npk

import (
	"bytes"
	"github.com/dof-wiki/godof/npk/img"
	"github.com/dof-wiki/godof/utils"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"github.com/dof-wiki/godof/utils/zlib_helper"
	"io"
)

type File struct {
	Name string

	f      io.ReadWriteSeeker
	offset int32
	size   int32
	data   []byte
}

func NewFile(f io.ReadWriteSeeker) (*File, error) {
	file := &File{
		f: f,
	}
	binary_helper.ReadAny(f, &file.offset)
	binary_helper.ReadAny(f, &file.size)
	buf := make([]byte, 256)
	f.Read(buf)
	nameData := decryptFileName(buf)
	nameData = utils.TrimBytesZeros(nameData)
	name, err := utils.Decode(nameData, "euc_kr")
	if err != nil {
		return nil, err
	} else {
		//idx := bytes.IndexByte([]byte(name), 0x00)
		//name = name[:idx]
	}
	name = utils.TrimStringZeros(name)
	file.Name = name
	return file, nil
}

func decryptFileName(data []byte) []byte {
	data = zlib_helper.FillBytes(data, 256)
	result := make([]byte, 256)
	for i := 0; i < 256; i++ {
		result[i] = data[i] ^ NPK_FILENAME_DECORD_FLAG[i]
	}
	return result
}

func (f *File) GetData() []byte {
	if f.data == nil {
		f.load()
	}
	return f.data
}

func (f *File) load() {
	f.f.Seek(int64(f.offset), io.SeekStart)
	f.data = make([]byte, f.size)
	f.f.Read(f.data)
}

func (f *File) ToIMG() (*img.Img, error) {
	reader := bytes.NewReader(f.GetData())
	return img.Open(reader)
}
