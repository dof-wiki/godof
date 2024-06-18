package image

import (
	"github.com/dof-wiki/godof/utils/binary_helper"
	"io"
)

type ZlibSpriteImage struct {
	c *CommonImage

	keep        int32
	spriteIndex int32
	left        int32
	top         int32
	right       int32
	bottom      int32
	rotate      int32
}

func (z *ZlibSpriteImage) GetData() []byte {
	return z.c.GetData()
}

func (z *ZlibSpriteImage) GetFormat() int32 {
	return z.c.GetFormat()
}

func (z *ZlibSpriteImage) WH() (int, int) {
	return z.c.WH()
}

func (z *ZlibSpriteImage) SetOffset(offset int64) {
	z.c.SetOffset(offset)
}

func (z *ZlibSpriteImage) GetSize() int32 {
	return z.c.GetSize()
}

func (z *ZlibSpriteImage) FixSize() {
	z.c.FixSize()
}

func (z *ZlibSpriteImage) GetSpriteIndex() int32 {
	return z.spriteIndex
}

func (z *ZlibSpriteImage) GetBox() [4]int32 {
	return [4]int32{
		z.left,
		z.top,
		z.right,
		z.bottom,
	}
}

func NewZlibSpriteImage(reader io.ReadSeeker, format int32) (Image, error) {
	c, err := NewCommonImage(reader, format)
	if err != nil {
		return nil, err
	}
	i := &ZlibSpriteImage{
		c: c.(*CommonImage),
	}
	binary_helper.ReadAny(reader, &i.keep)
	binary_helper.ReadAny(reader, &i.spriteIndex)
	binary_helper.ReadAny(reader, &i.left)
	binary_helper.ReadAny(reader, &i.top)
	binary_helper.ReadAny(reader, &i.right)
	binary_helper.ReadAny(reader, &i.bottom)
	binary_helper.ReadAny(reader, &i.rotate)
	return i, nil
}
