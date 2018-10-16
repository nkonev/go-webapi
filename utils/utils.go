package utils

import "golang.org/x/crypto/bcrypt"

type H map[string]interface{}

// TODO salt
func HashPassword(password string) (string, error) {
	passwordHash, passwordHashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if passwordHashErr != nil {
		return "", passwordHashErr
	}
	return string(passwordHash), nil
}
