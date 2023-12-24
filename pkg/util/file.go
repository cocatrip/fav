package util

import (
	"io/fs"
	"os"

	"github.com/cocatrip/fav/pkg/logger"
)

var log = logger.GetLogger()

func GetFileSize(file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		log.Errorf("cannot stat file: %v", err)
	}

	return int(fileStat.Size())
}

func GetFileMode(file *os.File) fs.FileMode {
	fileStat, err := file.Stat()
	if err != nil {
		log.Errorf("cannot stat file: %v", err)
	}

	return fileStat.Mode()
}
