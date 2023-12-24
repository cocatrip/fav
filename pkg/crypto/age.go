package crypto

import (
	"bytes"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/cocatrip/fav/pkg/logger"
	"github.com/cocatrip/fav/pkg/util"
)

var log = logger.GetLogger()

func ageGetIdentityFromFile(secretFile string) *age.X25519Identity {
	secretKey, err := os.Open(secretFile)
	if err != nil {
		return nil
	}

	identities, err := age.ParseIdentities(secretKey)
	if err != nil {
		return nil
	}

	if len(identities) == 0 {
		log.Errorln("no identities found in the input")
	}

	ageIdentity := identities[0]
	identity, ok := ageIdentity.(*age.X25519Identity)
	if !ok {
		log.Errorf("error converting age identity to x25519 identity")
	}

	return identity
}

func ageEncrypt(identity *age.X25519Identity, file *os.File) error {
	recipient := identity.Recipient()

	buf := &bytes.Buffer{}
	armorWriter := armor.NewWriter(buf)

	w, err := age.Encrypt(armorWriter, recipient)
	if err != nil {
		return err
	}

	content := make([]byte, util.GetFileSize(file))
	if _, err := file.Read(content); err != nil {
		return err
	}

	if _, err := w.Write(content); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	if err := armorWriter.Close(); err != nil {
		return err
	}

	encryptedFileName := file.Name() + ".age"
	encryptedFileMode := util.GetFileMode(file)
	if err := os.WriteFile(encryptedFileName, buf.Bytes(), encryptedFileMode); err != nil {
		return err
	}

	log.Infoln(file.Name())
	if err := os.Remove(file.Name()); err != nil {
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

	decryptedFileName := strings.TrimSuffix(file.Name(), ".age")
	decryptedFileMode := util.GetFileMode(file)
	if err := os.WriteFile(decryptedFileName, out.Bytes(), decryptedFileMode); err != nil {
		return err
	}

	if err := os.Remove(file.Name()); err != nil {
		return err
	}

	return nil
}
