package akevitt

import "golang.org/x/crypto/bcrypt"

// Hashes password using Bcrypt algorithm
func hashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// Compares hash and password
func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
