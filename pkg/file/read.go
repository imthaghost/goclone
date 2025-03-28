package file

import (
	"errors"
	"os"
)

// Exists checks if a file or directory exists
func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}

func GetFileContent(filepath string) string {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return ""
	}
	collectorBytes, _ := os.ReadFile(filepath)
	return string(collectorBytes[:])
}
