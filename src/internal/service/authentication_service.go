package service

import (
	"Notifications/src/internal/domain"
	"context"
	"errors"
	"uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

type RegisterRequest struct {
	Name        string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	DeviceToken string
}

type AuthenticationService struct {
	userRepo        UserRepository
	passwordEncoder PasswordEncoder
	tokenProvider   TokenProvider
}

func NewAuthenticationService(userRepo UserRepository, passwordEncoder PasswordEncoder, tokenProvider TokenProvider) *AuthenticationService {
	return &AuthenticationService{
		userRepo:        userRepo,
		passwordEncoder: passwordEncoder,
		tokenProvider:   tokenProvider,
	}
}

func (s *AuthenticationService) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	passwordHash, err := s.passwordEncoder.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		DeviceToken:  req.DeviceToken,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthenticationService) Authenticate(ctx context.Context, email, plaintextPassword string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !s.passwordEncoder.Verify(plaintextPassword, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", ErrInvalidCredentials
	}

	return s.tokenProvider.GenerateToken(*user)
}

func (s *AuthenticationService) ValidateToken(token string) (*TokenClaims, error) {
	return s.tokenProvider.VerifyToken(token)
}
