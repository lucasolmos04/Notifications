package service_test

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"Notifications/src/internal/service/mocks"
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticationService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	encoder := mocks.NewMockPasswordEncoder(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "lucas@example.com").Return(nil, service.ErrNotFound)
	encoder.EXPECT().Hash("supersecret").Return("hashed-supersecret", nil)

	var createdUser *domain.User
	userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *domain.User) error {
			createdUser = u
			return nil
		})

	svc := service.NewAuthenticationService(userRepo, encoder, mocks.NewMockTokenProvider(ctrl))

	user, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "Lucas", LastName: "Olmos", Email: "lucas@example.com", Password: "supersecret",
	})

	require.NoError(t, err)
	assert.Equal(t, "hashed-supersecret", user.PasswordHash)
	assert.True(t, user.IsActive, "a newly registered user should be active")
	require.NotNil(t, createdUser, "userRepo.Create should be called with the new user")
	assert.Equal(t, "lucas@example.com", createdUser.Email)
}

func TestAuthenticationService_Register_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "taken@example.com").
		Return(&domain.User{Email: "taken@example.com"}, nil)
	// userRepo.Create is intentionally never EXPECTed: gomock fails the test if it's called.

	svc := service.NewAuthenticationService(userRepo, mocks.NewMockPasswordEncoder(ctrl), mocks.NewMockTokenProvider(ctrl))

	_, err := svc.Register(context.Background(), service.RegisterRequest{Email: "taken@example.com", Password: "x"})

	require.ErrorIs(t, err, service.ErrEmailAlreadyExists)
}

func TestAuthenticationService_Register_HashFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	encoder := mocks.NewMockPasswordEncoder(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "a@b.com").Return(nil, service.ErrNotFound)
	encoder.EXPECT().Hash("x").Return("", errors.New("bcrypt broke"))

	svc := service.NewAuthenticationService(userRepo, encoder, mocks.NewMockTokenProvider(ctrl))

	_, err := svc.Register(context.Background(), service.RegisterRequest{Email: "a@b.com", Password: "x"})

	require.Error(t, err)
}

func TestAuthenticationService_Authenticate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	encoder := mocks.NewMockPasswordEncoder(ctrl)
	tokenProvider := mocks.NewMockTokenProvider(ctrl)

	user := domain.User{ID: uuid.New(), Email: "lucas@example.com", PasswordHash: "hashed", IsActive: true}

	userRepo.EXPECT().FindByEmail(gomock.Any(), "lucas@example.com").Return(&user, nil)
	encoder.EXPECT().Verify("correct", "hashed").Return(true)
	tokenProvider.EXPECT().GenerateToken(user).Return("signed-token", nil)

	svc := service.NewAuthenticationService(userRepo, encoder, tokenProvider)

	token, err := svc.Authenticate(context.Background(), "lucas@example.com", "correct")

	require.NoError(t, err)
	assert.Equal(t, "signed-token", token)
}

func TestAuthenticationService_Authenticate_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)

	userRepo.EXPECT().FindByEmail(gomock.Any(), "ghost@example.com").Return(nil, service.ErrNotFound)

	svc := service.NewAuthenticationService(userRepo, mocks.NewMockPasswordEncoder(ctrl), mocks.NewMockTokenProvider(ctrl))

	_, err := svc.Authenticate(context.Background(), "ghost@example.com", "whatever")

	require.ErrorIs(t, err, service.ErrInvalidCredentials, "should not leak a 'user not found'-specific error")
}

func TestAuthenticationService_Authenticate_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	encoder := mocks.NewMockPasswordEncoder(ctrl)

	user := domain.User{Email: "lucas@example.com", PasswordHash: "hashed", IsActive: true}
	userRepo.EXPECT().FindByEmail(gomock.Any(), "lucas@example.com").Return(&user, nil)
	encoder.EXPECT().Verify("wrong", "hashed").Return(false)

	svc := service.NewAuthenticationService(userRepo, encoder, mocks.NewMockTokenProvider(ctrl))

	_, err := svc.Authenticate(context.Background(), "lucas@example.com", "wrong")

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestAuthenticationService_Authenticate_InactiveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	encoder := mocks.NewMockPasswordEncoder(ctrl)

	user := domain.User{Email: "lucas@example.com", PasswordHash: "hashed", IsActive: false}
	userRepo.EXPECT().FindByEmail(gomock.Any(), "lucas@example.com").Return(&user, nil)
	encoder.EXPECT().Verify("correct", "hashed").Return(true)

	svc := service.NewAuthenticationService(userRepo, encoder, mocks.NewMockTokenProvider(ctrl))

	_, err := svc.Authenticate(context.Background(), "lucas@example.com", "correct")

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestAuthenticationService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenProvider := mocks.NewMockTokenProvider(ctrl)

	claims := &service.TokenClaims{UserID: uuid.New()}
	tokenProvider.EXPECT().VerifyToken("abc").Return(claims, nil)

	svc := service.NewAuthenticationService(mocks.NewMockUserRepository(ctrl), mocks.NewMockPasswordEncoder(ctrl), tokenProvider)

	result, err := svc.ValidateToken("abc")

	require.NoError(t, err)
	assert.Equal(t, claims.UserID, result.UserID)
}
