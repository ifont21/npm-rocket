package domain

type FileRepository interface {
	GetPlainFileContent(libPath string) ([]byte, error)
	GetJSONFileContent(libPath string) (map[string]interface{}, error)
	UpdateTopOfTheFileContent(newContent string, libPath string) error
}
