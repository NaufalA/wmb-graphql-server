package util

import "golang.org/x/crypto/bcrypt"

type PasswordUtil struct {}

func (*PasswordUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (*PasswordUtil) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}