package auth

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"fmt"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTTokenProvider struct {
	secretKey []byte
	tokenTTL  time.Duration
}

func NewJWTTokenProvider(secretKey string, tokenTTL time.Duration) *JWTTokenProvider {
	return &JWTTokenProvider{
		secretKey: []byte(secretKey),
		tokenTTL:  tokenTTL,
	}
}

func (p *JWTTokenProvider) GenerateToken(user domain.User) (string, error) {
	now := time.Now()
	claims := jwtClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secretKey)
}

func (p *JWTTokenProvider) VerifyToken(tokenString string) (*service.TokenClaims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return p.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &service.TokenClaims{UserID: userID}, nil
}
