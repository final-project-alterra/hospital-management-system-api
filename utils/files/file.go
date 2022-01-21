package files

import (
	"os"
)

func DoesExist(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}
