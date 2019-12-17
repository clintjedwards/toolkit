package password

import "golang.org/x/crypto/bcrypt"

// CheckPasswordHash validates a password against the stored hash
// to verify the user is authorized
func CheckPasswordHash(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

// HashPassword converts a byte string password into a bcrypt hash
// which is then stored as the only form of password
func HashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	return hash, err
}
