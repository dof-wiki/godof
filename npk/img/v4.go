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
	v2         *ImgV2
	colorBoard *ColorBoard
}

func newImgV4() *ImgV4 {
	return &ImgV4{
		v2: new(ImgV2),
	}
}

func (i *ImgV4) loadColorBoard(img *Img) {

}

func (i *ImgV4) onOpen(img *Img) error {
	i.colorBoard = newColorBoard(img.f)
	if err := i.v2.onOpen(img); err != nil {
		return err
	}
	return nil
}

func (i *ImgV4) build(i2 image.Image) ([]byte, int, int, error) {
	_, isZlib := i2.(*image.ZlibImage)
	var raw []byte
	var err error
	if isZlib && len(i.colorBoard.color) > 0 {
		raw, err = formatter.FormatToRawIndexes(i2.GetData(), i.colorBoard.color)
	} else {
		raw, err = formatter.FormatToRaw(i2.GetData(), i2.GetFormat())
	}
	if err != nil {
		return nil, 0, 0, nil
	}
	w, h := i2.WH()
	return raw, w, h, nil
}
