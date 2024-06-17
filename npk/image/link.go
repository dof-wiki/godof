package image

import (
	image2 "image"
	"io"
)

type LinkImage struct {
}

func (l *LinkImage) GetData() []byte {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) GetFormat() int32 {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) WH() (int, int) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) Build() (image2.Image, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) SetOffset(offset int64) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) GetSize() int32 {
	//TODO implement me
	panic("implement me")
}

func (l *LinkImage) FixSize() {
	//TODO implement me
	panic("implement me")
}

func NewLinkImage(reader io.ReadSeeker) (Image, error) {
	return &LinkImage{}, nil
}

func (l *LinkImage) LoadLink() {

}
