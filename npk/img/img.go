package img

import (
	"errors"
	"fmt"
	"github.com/dof-wiki/godof/npk/image"
	"github.com/dof-wiki/godof/utils"
	"github.com/dof-wiki/godof/utils/binary_helper"
	image2 "image"
	"io"
)

const (
	IMG_MAGIC_OLD = "Neople Image File"
	IMG_MAGIC     = "Neople Img File"
)

type ImgIO interface {
	onOpen(*Img) error
	build(image.Image) ([]byte, int, int, error)
}

func newImgIO(version int32) ImgIO {
	switch version {
	case 1:
		return new(ImgV1)
	case 2:
		return new(ImgV2)
	case 4:
		return newImgV4()
	}
	return nil
}

type Img struct {
	f          io.ReadSeeker
	version    int32
	keep       int32
	imagesSize int32
	imageCount int32
	io         ImgIO

	Images []image.Image
}

func (i *Img) Build(i2 image.Image) (image2.Image, error) {
	raw, w, h, err := i.io.build(i2)
	if err != nil {
		return nil, err
	}
	i3 := image2.NewRGBA(image2.Rect(0, 0, w, h))
	copy(i3.Pix, raw)
	return i3, nil
}

func newImg(f io.ReadSeeker, version, keep, imagesSize, imageCount int32) (*Img, error) {
	imgIO := newImgIO(version)
	if imgIO == nil {
		return nil, errors.New(fmt.Sprintf("IMG version %d not support.", version))
	}
	img := &Img{
		f:          f,
		version:    version,
		keep:       keep,
		imagesSize: imagesSize,
		imageCount: imageCount,
		io:         imgIO,
		Images:     make([]image.Image, 0),
	}
	if err := img.io.onOpen(img); err != nil {
		return nil, err
	}
	return img, nil
}

func Open(f io.ReadSeeker) (*Img, error) {
	magic := binary_helper.ReadStringByLen(f, 16)
	magic = utils.TrimStringZeros(magic)
	if magic != IMG_MAGIC && magic != IMG_MAGIC_OLD[:16] {
		return nil, errors.New("IMG header err")
	}

	var imagesSize int32
	if magic == IMG_MAGIC {
		binary_helper.ReadAny(f, &imagesSize)
	} else {
		binary_helper.ReadStringByLen(f, 2)
		var unknown int16
		binary_helper.ReadAny(f, &unknown)
	}

	var keep, version, imageCount int32
	binary_helper.ReadAny(f, &keep)
	binary_helper.ReadAny(f, &version)
	binary_helper.ReadAny(f, &imageCount)
	return newImg(f, version, keep, imagesSize, imageCount)
}
