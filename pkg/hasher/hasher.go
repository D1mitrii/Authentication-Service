package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	salt int
}

func New(salt int) *Hasher {
	return &Hasher{
		salt: salt,
	}
}

// Hash generate hash from password using bcrypt
func (h *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.salt)
	return string(hash), err
}

// Compare plain password with hashed password
func (h *Hasher) Compare(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
