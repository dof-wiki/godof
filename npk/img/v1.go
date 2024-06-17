package img

import (
	"github.com/dof-wiki/godof/npk/image"
	"io"
	"log"
)

type ImgV1 struct {
}

func (i *ImgV1) onOpen(img *Img) error {
	images := make([]image.Image, 0, img.imageCount)
	for j := int32(0); j < img.imageCount; j++ {
		newImage, err := image.NewImage(img.f)
		if err != nil {
			log.Printf("open image err %v", err)
			continue
		}
		offset, _ := img.f.Seek(0, io.SeekCurrent)
		newImage.SetOffset(offset)
		img.f.Seek(int64(newImage.GetSize()), io.SeekCurrent)
		images = append(images, newImage)

		if link, ok := newImage.(*image.LinkImage); ok {
			link.LoadLink()
		}
	}

	img.Images = images
	return nil
}
