package img

import (
	"github.com/dof-wiki/godof/npk/image"
	"github.com/dof-wiki/godof/npk/image/formatter"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"github.com/dof-wiki/godof/utils/zlib_helper"
	"github.com/samber/lo"
	"io"
)

type Sprite struct {
	keep     int32
	fmt      int32
	index    int32
	dataSize int32
	rawSize  int32
	w        int32
	h        int32
	offset   int64
	data     []byte
	zipData  []byte
	reader   io.ReadSeeker
}

func (s *Sprite) GetData() []byte {
	if s.data == nil {
		s.loadData()
	}
	return s.data
}

func (s *Sprite) loadData() {
	s.reader.Seek(s.offset, io.SeekStart)
	s.zipData = make([]byte, s.dataSize)
	s.reader.Read(s.zipData)
	s.data, _ = zlib_helper.Decompress(s.zipData)
}

func newSprite(reader io.ReadSeeker) *Sprite {
	s := &Sprite{
		reader: reader,
	}
	binary_helper.ReadAny(reader, &s.keep)
	binary_helper.ReadAny(reader, &s.fmt)
	binary_helper.ReadAny(reader, &s.index)
	binary_helper.ReadAny(reader, &s.dataSize)
	binary_helper.ReadAny(reader, &s.rawSize)
	binary_helper.ReadAny(reader, &s.w)
	binary_helper.ReadAny(reader, &s.h)
	return s
}

type ImgV5 struct {
	sprites    []*Sprite
	colorBoard *ColorBoard
}

func newImgV5() *ImgV5 {
	return &ImgV5{
		sprites: make([]*Sprite, 0),
	}
}

func (i *ImgV5) onOpen(img *Img) error {
	var spriteCount, fileSize int32
	binary_helper.ReadAny(img.f, &spriteCount)
	binary_helper.ReadAny(img.f, &fileSize)
	i.colorBoard = newColorBoard(img.f)

	for j := 0; j < int(spriteCount); j++ {
		i.sprites = append(i.sprites, newSprite(img.f))
	}

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
	for _, s := range i.sprites {
		s.offset = offset
		offset += int64(s.dataSize)
	}
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

func (i *ImgV5) build(i2 image.Image) ([]byte, int, int, string, error) {
	var raw []byte
	var err error
	if i3, ok := i2.(*image.ZlibSpriteImage); ok {
		format := ""
		sprite := i.sprites[i3.GetSpriteIndex()]
		if lo.Contains(image.IMAGE_FORMATS_DDS, int(sprite.fmt)) {
			format = "dds"
			raw = sprite.GetData()
		} else {
			box := i3.GetBox()
			raw, err = formatter.FormatToRawCrop(sprite.GetData(), sprite.fmt, sprite.w, box[0], box[1], box[2], box[3])
		}
		if err != nil {
			return nil, 0, 0, "", nil
		}
		w, h := i2.WH()
		return raw, w, h, format, nil
	}

	_, isZlib := i2.(*image.ZlibImage)
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
