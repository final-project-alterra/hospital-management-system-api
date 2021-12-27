package hash

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	const op errors.Op = "hash.Hash"
	const COST = 14

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), COST)
	if err != nil {
		return "", errors.E(op, errors.KindUnprocessable, err)
	}
	return string(hashed), nil
}

func Validate(hashed string, original string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(original))
	return err == nil
}
