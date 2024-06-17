package image

import (
	"errors"
	"fmt"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"github.com/samber/lo"
	image2 "image"
	"io"
)

type Image interface {
	SetOffset(offset int64)
	GetSize() int32
	FixSize()
	Build() (image2.Image, error)
}

var ImageInstanceMap = map[int32]func(reader io.ReadSeeker, format int32) (Image, error){
	IMAGE_EXTRA_NONE: NewCommonImage,
	IMAGE_EXTRA_ZLIB: NewZlibImage,
}

func NewImage(reader io.ReadSeeker) (Image, error) {
	var format int32
	binary_helper.ReadAny(reader, &format)

	if !lo.Contains(IMAGE_FORMATS_ALL, int(format)) {
		return nil, errors.New(fmt.Sprintf("Invalid Image format %d", format))
	}
	if format == IMAGE_FORMAT_LINK {
		// TODO
		return nil, errors.New("Not implemented")
	} else {
		var extra int32
		binary_helper.ReadAny(reader, &extra)
		if _, ok := ImageInstanceMap[extra]; !ok {
			return nil, errors.New(fmt.Sprintf("Invalid Image extra %d", extra))
		}
		return ImageInstanceMap[extra](reader, format)
	}
}
