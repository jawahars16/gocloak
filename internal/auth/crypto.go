package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type crypto struct{}

func NewCrypto() crypto {
	return crypto{}
}

func (c *crypto) GenerateFromPassword(password string) (string, error) {
	hashData, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return password, err
	}
	return string(hashData), nil
}

func (c *crypto) CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
