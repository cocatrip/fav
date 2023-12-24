package crypto

func Encrypt(secretFile string, file string) {
	identity := ageGetIdentityFromFile(secretFile)

	err := ageEncrypt(identity, file)
	if err != nil {
		return
	}
}
