package npk

import (
	"errors"
	"fmt"
	"github.com/dof-wiki/godof/utils"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"io"
)

type Npk struct {
	FileCount int32
	Files     []*File
}

func Open(f io.ReadWriteSeeker) (*Npk, error) {
	magic := binary_helper.ReadStringByLen(f, 16)
	magic = utils.TrimStringZeros(magic)
	if magic != NPK_MAGIC {
		fmt.Println(magic, len(magic), magic[len(magic)-1])
		fmt.Println(NPK_MAGIC, len(NPK_MAGIC))
		return nil, errors.New("file is not a valid npk")
	}
	n := &Npk{
		Files: make([]*File, 0),
	}
	binary_helper.ReadAny(f, &n.FileCount)

	for i := int32(0); i < n.FileCount; i++ {
		file, err := NewFile(f)
		if err != nil {
			return nil, err
		}
		n.Files = append(n.Files, file)
	}
	return n, nil
}
