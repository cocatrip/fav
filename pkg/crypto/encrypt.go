package crypto

import (
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/cocatrip/fav/pkg/config"
	"github.com/cocatrip/fav/pkg/logger"
	"github.com/cocatrip/fav/pkg/util"
	"github.com/cocatrip/fav/shamir"
	"github.com/spf13/viper"
)

/*
Encrypt will generate an age  identity  file  and  split  those  identity  using
shamir secret sharing while using the identity to encrypt  the  named  file.  If
successful then the splitted identity will be uploaded to buckets listed in  the
fav config file
*/
func Encrypt(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	identity, err := ageEncrypt(file)
	if err != nil {
		return err
	}

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	secretFiles, err := shamirSplit(identity, util.GetEncFileName(file), len(config.Storages))
	if err != nil {
		return err
	}

	for i, storage := range config.Storages {
		if err := storage.Upload(secretFiles[i]); err != nil {
			return err
		}
	}

	return err
}

func ageEncrypt(file *os.File) (*age.X25519Identity, error) {
	log := logger.GetLogger()

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}

	recipient := identity.Recipient()

	log.Infof(
		"running: age --encrypt --armor --recipient %s > %s && rm %s",
		recipient,
		util.GetEncFileName(file),
		file.Name(),
	)

	encFile, err := os.Create(util.GetEncFileName(file))
	if err != nil {
		return nil, err
	}

	armorWriter := armor.NewWriter(encFile)

	w, err := age.Encrypt(armorWriter, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(w, file); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if err := armorWriter.Close(); err != nil {
		return nil, err
	}

	if err := os.Remove(file.Name()); err != nil {
		return nil, err
	}

	return identity, nil
}

/*
shamirSplit acts  as  a  wrapper  for  splitting  an  age  identity  into  files
using the original file name as a suffix and  then  returns  an  array  of  file
name that the identity was splitted
*/
func shamirSplit(identity *age.X25519Identity, fileName string, parts int) ([]string, error) {
	sharedSecret, err := shamir.Split([]byte(identity.String()), parts, parts)
	if err != nil {
		return nil, err
	}

	secretFiles := make([]string, 0)

	for _, secret := range sharedSecret {
		secretFileName := util.GenerateSecretFileName(fileName)

		secretFile, err := os.Create(secretFileName)
		if err != nil {
			return nil, err
		}
		defer secretFile.Close()

		if _, err := secretFile.Write(secret); err != nil {
			return nil, err
		}

		secretFiles = append(secretFiles, secretFileName)
	}

	return secretFiles, nil
}
