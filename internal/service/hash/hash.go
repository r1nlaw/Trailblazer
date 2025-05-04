package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) bool
}

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func (b *BcryptHasher) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
