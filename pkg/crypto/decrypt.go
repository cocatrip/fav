package crypto

import "os"

func Decrypt(secretFile string, fileName string) error {
	identity := ageGetIdentityFromFile(secretFile)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := ageDecrypt(identity, file); err != nil {
		return err
	}

	return nil
}
