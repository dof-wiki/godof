package pvf_source

import (
	"os"
	"path"
)

type FileSystemSource struct {
	rootDir string
}

func NewFileSystemSource(rootDir string) *FileSystemSource {
	return &FileSystemSource{
		rootDir: rootDir,
	}
}

func (p *FileSystemSource) GetFileContent(filepath string) (string, error) {
	buf, err := os.ReadFile(path.Join(p.rootDir, filepath))
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (p *FileSystemSource) SaveFileContent(filepath, content string) error {
	f, err := os.Create(path.Join(p.rootDir, filepath))
	if err != nil {
		return nil
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
