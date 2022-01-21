package project

import (
	"os"
	"path/filepath"
)

var mainDir string

func loadMainDirectory() {
	var executableFilePath string
	var err error

	if mainDir != "" {
		return
	}

	if executableFilePath, err = os.Executable(); err != nil {
		panic(err.Error())
	}

	mainDir = filepath.Dir(executableFilePath)
}

var GetMainDir = func() string {
	if mainDir == "" {
		loadMainDirectory()
	}
	return mainDir
}
