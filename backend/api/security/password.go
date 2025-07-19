package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	return string(hash), err
}

func VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
