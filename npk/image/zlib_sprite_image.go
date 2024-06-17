package image

import (
	image2 "image"
	"io"
)

type ZlibSpriteImage struct {
}

func (z *ZlibSpriteImage) Build() (image2.Image, error) {
	//TODO implement me
	panic("implement me")
}

func (z *ZlibSpriteImage) SetOffset(offset int64) {
	//TODO implement me
	panic("implement me")
}

func (z *ZlibSpriteImage) GetSize() int32 {
	//TODO implement me
	panic("implement me")
}

func (z *ZlibSpriteImage) FixSize() {
	//TODO implement me
	panic("implement me")
}

func NewZlibSpriteImage(reader io.ReadSeeker) (Image, error) {
	return &ZlibSpriteImage{}, nil
}
