package helpers

import "golang.org/x/crypto/bcrypt"

func GenerateHashedPassword(password []byte) (string, error) {
	salt := 8

	hash, err := bcrypt.GenerateFromPassword(password, salt)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePass(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
