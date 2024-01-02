package crypto

import (
	"bytes"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/cocatrip/fav/pkg/config"
	"github.com/cocatrip/fav/pkg/util"
	"github.com/cocatrip/fav/shamir"
	"github.com/spf13/viper"
)

func Decrypt(name string) error {
	// read file
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download from buckets
	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	// Get md5sum from encrypted file
	md5Sum := util.GetMd5Sum(name)

	for _, storage := range config.Storages {
		if err := storage.Download(md5Sum); err != nil {
			return err
		}
	}

	secretFiles, err := util.GetFilesInDirectory(".tmp")
	if err != nil {
		return err
	}

	identity, err := shamirCombine(secretFiles)
	if err != nil {
		return err
	}

	if err := ageDecrypt(identity, file); err != nil {
		return err
	}

	if err := os.RemoveAll(".tmp"); err != nil {
		return err
	}

	return nil
}

func ageDecrypt(identity *age.X25519Identity, file *os.File) error {
	content := make([]byte, util.GetFileSize(file))
	if _, err := file.Read(content); err != nil {
		return err
	}

	out := &bytes.Buffer{}
	f := strings.NewReader(string(content))
	armorReader := armor.NewReader(f)

	r, err := age.Decrypt(armorReader, identity)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, r); err != nil {
		return err
	}

	if err := os.WriteFile(util.GetDecFileName(file), out.Bytes(), util.GetFileMode(file)); err != nil {
		return err
	}

	if err := os.Remove(file.Name()); err != nil {
		return err
	}

	return nil
}

/*
shamirCombine  acts  as  a  wrapper  for  combining  files  containing  splitted
identity into a single identity. it will loop through the  array  of  file  name
and  then  get  the  bytes  and   combine   them   and   parse   that   identity
*/
func shamirCombine(fileList []string) (*age.X25519Identity, error) {
	var sharedSecret [][]byte

	for _, fileName := range fileList {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		sharedSecret = append(sharedSecret, fileData)
	}

	identityBytes, err := shamir.Combine(sharedSecret)
	if err != nil {
		return nil, err
	}

	identity, err := age.ParseX25519Identity(string(identityBytes))
	if err != nil {
		return nil, err
	}

	return identity, nil
}
