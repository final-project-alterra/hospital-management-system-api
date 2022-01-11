package project

import (
	"os"
	"path/filepath"
)

func MainDirectory() (string, error) {
	var executableFilePath string
	var err error

	if executableFilePath, err = os.Executable(); err != nil {
		return "", err
	}

	return filepath.Dir(executableFilePath), nil
}
