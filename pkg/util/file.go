package util

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cocatrip/fav/pkg/logger"
	"k8s.io/apimachinery/pkg/util/rand"
)

var log = logger.GetLogger()

func GetFileSize(file *os.File) int {
	fileStat, err := file.Stat()
	if err != nil {
		log.Error(err)
	}

	return int(fileStat.Size())
}

func GetFileMode(file *os.File) fs.FileMode {
	fileStat, err := file.Stat()
	if err != nil {
		log.Error(err)
	}

	return fileStat.Mode()
}

func GetEncFileName(file *os.File) string {
	name := fmt.Sprintf("%s.age", file.Name())

	return name
}

func GetDecFileName(file *os.File) string {
	name := strings.TrimSuffix(file.Name(), ".age")

	return name
}

func GenerateSecretFileName(name string) string {
	f, err := os.ReadFile(name)
	if err != nil {
		log.Error(err)
	}

	md5Name := md5.Sum(f)
	randName := rand.String(len(name))

	secretFileName := fmt.Sprintf("%x.%s", md5Name, randName)

	return secretFileName
}

func GetMd5Sum(name string) string {
	f, err := os.ReadFile(name)
	if err != nil {
		log.Error(err)
	}

	md5Name := md5.Sum(f)

	return fmt.Sprintf("%x", md5Name)
}

func GetFilesInDirectory(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
