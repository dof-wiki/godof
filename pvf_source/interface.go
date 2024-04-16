package pvf_source

type PvfSource interface {
	GetFileContent(path string) (string, error)
	SaveFileContent(path, content string) error
}
