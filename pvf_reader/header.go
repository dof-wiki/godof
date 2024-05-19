package pvf_reader

import (
	"bytes"
	"fmt"
	"github.com/dof-wiki/godof/utils"
	"github.com/dof-wiki/godof/utils/binary_helper"
	"io"
	"os"
	"strings"
	"time"
)

type PvfHeader struct {
	uuidLen         int32
	uuid            string
	version         int32
	dirTreeLen      int32
	dirTreeCrc32    uint32
	fileCount       int32
	data            []byte
	FileTree        map[string]*FileTreeNode
	ContentStartIdx int64
}

func newPvfHeader() *PvfHeader {
	return &PvfHeader{
		FileTree: make(map[string]*FileTreeNode),
	}
}

func (h *PvfHeader) read(f *os.File) error {
	h.uuid, h.uuidLen = binary_helper.ReadStr(f)
	binary_helper.ReadAny(f, &h.version)
	binary_helper.ReadAny(f, &h.dirTreeLen)
	binary_helper.ReadAny(f, &h.dirTreeCrc32)
	binary_helper.ReadAny(f, &h.fileCount)

	fmt.Printf("uuidLen(%d), uuid(%s), version(%d), treeLen(%d), crc32(%d), fileCount(%d)\n", h.uuidLen, h.uuid, h.version, h.dirTreeLen, h.dirTreeCrc32, h.fileCount)

	buf := make([]byte, h.dirTreeLen)
	if _, err := f.Read(buf); err != nil {
		return err
	}
	h.data = binary_helper.DecryptCrc(buf, h.dirTreeCrc32)
	h.ContentStartIdx, _ = f.Seek(0, io.SeekCurrent)
	return h.loadFileTree()
}

func (h *PvfHeader) loadFileTree() error {
	st := time.Now()
	defer func() {
		fmt.Printf("load file tree cost %s\n", time.Since(st))
	}()
	buffer := bytes.NewBuffer(h.data)
	for i := 0; i < int(h.fileCount); i++ {
		var fn uint32
		binary_helper.ReadAny(buffer, &fn)
		var filePathLen uint32
		filePathBytes := binary_helper.ReadBytesByLen(buffer, &filePathLen)
		filePath, err := utils.DecodeCP949(filePathBytes)
		if err != nil {
			return err
		}
		filePath, _ = strings.CutPrefix(strings.ToLower(filePath), "/")
		var fileLength uint32
		binary_helper.ReadAny(buffer, &fileLength)
		fileLength = (fileLength + 3) & 0xFFFFFFFC
		var fileCrc32 uint32
		binary_helper.ReadAny(buffer, &fileCrc32)
		var relativeOffset uint32
		binary_helper.ReadAny(buffer, &relativeOffset)
		node := &FileTreeNode{
			Fn:             fn,
			Path:           filePath,
			Length:         fileLength,
			Crc32:          fileCrc32,
			RelativeOffset: relativeOffset,
		}
		h.FileTree[filePath] = node
	}
	return nil
}
