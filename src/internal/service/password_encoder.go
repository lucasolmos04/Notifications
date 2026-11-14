package service

//go:generate mockgen -source=password_encoder.go -destination=mocks/password_encoder_mock.go -package=mocks

type PasswordEncoder interface {
	Hash(plaintext string) (string, error)
	Verify(plaintext, hash string) bool
}
