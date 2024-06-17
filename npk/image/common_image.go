package image

import (
	"github.com/dof-wiki/godof/npk/image/formatter"
	"github.com/dof-wiki/godof/utils/binary_helper"
	image2 "image"
	"io"
)

type CommonImage struct {
	format int32
	reader io.ReadSeeker
	offset int64
	size   int32
	data   []byte

	w  int32
	h  int32
	x  int32
	y  int32
	mw int32
	mh int32
}

func (c *CommonImage) loadData() {
	c.reader.Seek(c.offset, io.SeekStart)
	c.data = make([]byte, c.size)
	c.reader.Read(c.data)
}

func (c *CommonImage) GetData() []byte {
	if c.data == nil {
		c.loadData()
	}
	return c.data
}

func (c *CommonImage) GetFormat() int32 {
	return c.format
}

func (c *CommonImage) Build() (image2.Image, error) {
	data := c.GetData()
	raw, err := formatter.FormatToRaw(data, c.format)
	if err != nil {
		return nil, err
	}
	i := image2.NewRGBA(image2.Rect(0, 0, int(c.w), int(c.h)))
	copy(i.Pix, raw)
	return i, nil
}

func (c *CommonImage) WH() (int, int) {
	return int(c.w), int(c.h)
}

func (c *CommonImage) SetOffset(offset int64) {
	c.offset = offset
}

func (c *CommonImage) GetSize() int32 {
	return c.size
}

func (c *CommonImage) FixSize() {
	c.size = c.w * c.h * PIX_SIZE[c.format]
}

func NewCommonImage(reader io.ReadSeeker, format int32) (Image, error) {
	i := new(CommonImage)
	i.format = format
	i.reader = reader
	binary_helper.ReadAny(reader, &i.w)
	binary_helper.ReadAny(reader, &i.h)
	binary_helper.ReadAny(reader, &i.size)
	binary_helper.ReadAny(reader, &i.x)
	binary_helper.ReadAny(reader, &i.y)
	binary_helper.ReadAny(reader, &i.mw)
	binary_helper.ReadAny(reader, &i.mh)
	return i, nil
}
