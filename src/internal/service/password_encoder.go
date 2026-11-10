package service

type PasswordEncoder interface {
	Hash(plaintext string) (string, error)
	Verify(plaintext, hash string) bool
}
