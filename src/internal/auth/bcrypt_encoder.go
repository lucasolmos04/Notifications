package auth

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordEncoder struct{}

func NewBcryptPasswordEncoder() *BcryptPasswordEncoder {
	return &BcryptPasswordEncoder{}
}

func (e *BcryptPasswordEncoder) Hash(plaintext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (e *BcryptPasswordEncoder) Verify(plaintext, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
