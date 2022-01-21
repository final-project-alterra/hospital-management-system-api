package files

import (
	"os"
)

var DoesExist = func(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

var Remove = func(path string) error {
	return os.Remove(path)
}
