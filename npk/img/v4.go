package img

import (
	"github.com/dof-wiki/godof/npk/image"
	"github.com/dof-wiki/godof/npk/image/formatter"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"io"
)

type ColorBoard struct {
	color [][]uint8
}

func newColorBoard(reader io.ReadSeeker) *ColorBoard {
	c := &ColorBoard{
		color: make([][]uint8, 0),
	}
	var count int32
	binary_helper.ReadAny(reader, &count)
	for i := 0; i < int(count); i++ {
		var r, g, b, a uint8
		binary_helper.ReadAny(reader, &r)
		binary_helper.ReadAny(reader, &g)
		binary_helper.ReadAny(reader, &b)
		binary_helper.ReadAny(reader, &a)
		c.color = append(c.color, []uint8{r, g, b, a})
	}
	return c
}

type ImgV4 struct {
	colorBoard *ColorBoard
}

func newImgV4() *ImgV4 {
	return &ImgV4{}
}

func (i *ImgV4) loadColorBoard(img *Img) {

}

func (i *ImgV4) onOpen(img *Img) error {
	i.colorBoard = newColorBoard(img.f)
	images := make([]image.Image, 0, img.imageCount)
	for j := int32(0); j < img.imageCount; j++ {
		newImage, err := image.NewImage(img.f)
		if err != nil {
			return err
		}
		newImage.FixSize()
		images = append(images, newImage)
	}
	img.Images = images

	offset, _ := img.f.Seek(0, io.SeekCurrent)
	for _, item := range img.Images {
		if _, ok := item.(*image.LinkImage); ok {
			continue
		}
		if _, ok := item.(*image.ZlibSpriteImage); ok {
			continue
		}
		item.SetOffset(offset)
		offset += int64(item.GetSize())
	}
	return nil
}

func (i *ImgV4) build(i2 image.Image) ([]byte, int, int, string, error) {
	_, isZlib := i2.(*image.ZlibImage)
	var raw []byte
	var err error
	if isZlib && len(i.colorBoard.color) > 0 {
		raw, err = formatter.FormatToRawIndexes(i2.GetData(), i.colorBoard.color)
	} else {
		raw, err = formatter.FormatToRaw(i2.GetData(), i2.GetFormat())
	}
	if err != nil {
		return nil, 0, 0, "", nil
	}
	w, h := i2.WH()
	return raw, w, h, "", nil
}
