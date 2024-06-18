package img

import (
	"github.com/dof-wiki/godof/npk/image"
	"github.com/dof-wiki/godof/npk/image/formatter"
)

type ImgV2 struct {
}

func (i *ImgV2) onOpen(img *Img) error {
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

	offset := img.imagesSize + 32
	for _, item := range img.Images {
		if _, ok := item.(*image.LinkImage); ok {
			continue
		}
		if _, ok := item.(*image.ZlibSpriteImage); ok {
			continue
		}
		item.SetOffset(int64(offset))
		offset += item.GetSize()
	}
	return nil
}

func (i *ImgV2) build(i2 image.Image) ([]byte, int, int, string, error) {
	data := i2.GetData()
	raw, err := formatter.FormatToRaw(data, i2.GetFormat())
	if err != nil {
		return nil, 0, 0, "", nil
	}
	w, h := i2.WH()
	return raw, w, h, "", nil
}
