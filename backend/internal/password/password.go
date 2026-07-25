package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(plainText),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Verify(hash string, plainText string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plainText),
	)
}
