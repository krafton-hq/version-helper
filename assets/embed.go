package assets

import "embed"

// Embedded contains embedded templates
//go:embed *
var Embedded embed.FS

func GetFile(fileName string) (string, error) {
	buf, err := Embedded.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
