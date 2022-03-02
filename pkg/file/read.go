package file

import (
	"errors"
	"os"
)

func GetFileContent(filepath string) string {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return ""
	}
	collectorBytes, _ := os.ReadFile(filepath)
	return string(collectorBytes[:])
}
