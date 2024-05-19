package pvf_reader

type FileTreeNode struct {
	Fn             uint32
	Path           string
	Length         uint32
	Crc32          uint32
	RelativeOffset uint32
}
