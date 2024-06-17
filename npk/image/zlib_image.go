package image

import (
	"github.com/dof-wiki/godof/npk/image/formatter"
	"github.com/dof-wiki/godof/utils/zlib_helper"
	image2 "image"
	"io"
)

type ZlibImage struct {
	c *CommonImage

	zipData []byte
}

func (z *ZlibImage) SetOffset(offset int64) {
	z.c.SetOffset(offset)
}

func (z *ZlibImage) GetSize() int32 {
	return z.c.GetSize()
}

func (z *ZlibImage) FixSize() {
}

func (z *ZlibImage) GetData() []byte {
	if z.c.data == nil {
		z.loadData()
	}
	return z.c.data
}

func (z *ZlibImage) Build() (image2.Image, error) {
	data := z.GetData()
	raw, err := formatter.FormatToRaw(data, z.c.format)
	if err != nil {
		return nil, err
	}
	i := image2.NewRGBA(image2.Rect(0, 0, int(z.c.w), int(z.c.h)))
	copy(i.Pix, raw)
	return i, nil
}

func (z *ZlibImage) loadData() {
	z.c.loadData()
	z.zipData = z.c.data
	z.c.data, _ = zlib_helper.Decompress(z.zipData)
}

func NewZlibImage(reader io.ReadSeeker, format int32) (Image, error) {
	c, err := NewCommonImage(reader, format)
	if err != nil {
		return nil, err
	}
	return &ZlibImage{
		c: c.(*CommonImage),
	}, nil
}
