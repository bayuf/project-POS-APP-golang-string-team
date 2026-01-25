package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashString(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("fail hashing password")
	}

	return string(bytes), nil
}

func CheckString(hashed, pass string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err == nil
}
