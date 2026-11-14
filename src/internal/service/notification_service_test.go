package service_test

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"Notifications/src/internal/service/mocks"
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNotificationService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	sender := mocks.NewMockNotificationSender(ctrl)

	userID := uuid.New()
	user := domain.User{ID: userID, Email: "lucas@example.com"}

	userRepo.EXPECT().FindByID(gomock.Any(), userID).Return(&user, nil)
	sender.EXPECT().Validate(gomock.Any(), gomock.Any(), user).Return(nil)
	repo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Notification{})).
		DoAndReturn(func(ctx context.Context, n domain.Notification) error {
			assert.Equal(t, domain.Pending, n.Status, "should be persisted as Pending before sending")
			return nil
		})
	sender.EXPECT().Send(gomock.Any(), gomock.Any(), user).
		Return(&service.SenderResult{Channel: domain.ChannelEmail, SendAt: time.Now(), Detail: "<html></html>"}, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(domain.Notification{})).
		DoAndReturn(func(ctx context.Context, n domain.Notification) error {
			assert.Equal(t, domain.Sent, n.Status, "the post-send update should persist Sent")
			return nil
		})

	svc := service.NewNotificationService(repo, userRepo, service.SenderRegistry{domain.ChannelEmail: sender})

	result, err := svc.Create(context.Background(), service.CreateNotificationRequest{
		UserID:  userID,
		Title:   "Hola",
		Content: "Contenido",
		Channel: domain.ChannelEmail,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.Sent, result.Status)
}

func TestNotificationService_Create_SendFailureMarksStatusFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	sender := mocks.NewMockNotificationSender(ctrl)

	userID := uuid.New()
	user := domain.User{ID: userID}

	userRepo.EXPECT().FindByID(gomock.Any(), userID).Return(&user, nil)
	sender.EXPECT().Validate(gomock.Any(), gomock.Any(), user).Return(nil)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	sender.EXPECT().Send(gomock.Any(), gomock.Any(), user).Return(nil, errors.New("provider unreachable"))

	var persistedStatus domain.StatusNotification
	repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(domain.Notification{})).
		DoAndReturn(func(ctx context.Context, n domain.Notification) error {
			persistedStatus = n.Status
			return nil
		})

	svc := service.NewNotificationService(repo, userRepo, service.SenderRegistry{domain.ChannelEmail: sender})

	result, err := svc.Create(context.Background(), service.CreateNotificationRequest{
		UserID: userID, Title: "t", Content: "c", Channel: domain.ChannelEmail,
	})

	require.NoError(t, err, "should not fail the Create call even when the send attempt fails")
	assert.Equal(t, domain.Failed, result.Status)
	assert.Equal(t, domain.Failed, persistedStatus)
}

func TestNotificationService_Create_UnsupportedChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	// No .EXPECT() calls set on repo/userRepo: gomock fails the test if either
	// is called, which asserts we bail out before touching persistence.

	svc := service.NewNotificationService(repo, userRepo, service.SenderRegistry{})

	_, err := svc.Create(context.Background(), service.CreateNotificationRequest{
		UserID: uuid.New(), Channel: domain.Channel("fax"),
	})

	require.ErrorIs(t, err, service.ErrUnsupportedChannel)
}

func TestNotificationService_Create_ValidateFailsPreventsPersistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	sender := mocks.NewMockNotificationSender(ctrl)

	userID := uuid.New()
	user := domain.User{ID: userID}

	userRepo.EXPECT().FindByID(gomock.Any(), userID).Return(&user, nil)
	sender.EXPECT().Validate(gomock.Any(), gomock.Any(), user).Return(errors.New("invalid recipient"))
	// repo.Create is intentionally never EXPECTed: gomock fails the test if it's called.

	svc := service.NewNotificationService(repo, userRepo, service.SenderRegistry{domain.ChannelEmail: sender})

	_, err := svc.Create(context.Background(), service.CreateNotificationRequest{
		UserID: userID, Channel: domain.ChannelEmail,
	})

	require.Error(t, err)
}

func TestNotificationService_Update_OwnershipMismatchReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)

	ownerID := uuid.New()
	otherUserID := uuid.New()
	notificationID := uuid.New()

	repo.EXPECT().FindByID(gomock.Any(), notificationID).
		Return(&domain.Notification{ID: notificationID, UserID: ownerID, Channel: domain.ChannelEmail}, nil)

	svc := service.NewNotificationService(repo, mocks.NewMockUserRepository(ctrl), service.SenderRegistry{domain.ChannelEmail: mocks.NewMockNotificationSender(ctrl)})

	_, err := svc.Update(context.Background(), otherUserID, notificationID, service.UpdateNotificationRequest{
		Channel: domain.ChannelEmail,
	})

	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestNotificationService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)

	userID := uuid.New()
	notificationID := uuid.New()

	repo.EXPECT().FindByID(gomock.Any(), notificationID).
		Return(&domain.Notification{ID: notificationID, UserID: userID, Title: "old", Channel: domain.ChannelEmail}, nil)

	var persistedTitle string
	repo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(domain.Notification{})).
		DoAndReturn(func(ctx context.Context, n domain.Notification) error {
			persistedTitle = n.Title
			return nil
		})

	svc := service.NewNotificationService(repo, mocks.NewMockUserRepository(ctrl), service.SenderRegistry{domain.ChannelEmail: mocks.NewMockNotificationSender(ctrl)})

	result, err := svc.Update(context.Background(), userID, notificationID, service.UpdateNotificationRequest{
		Title: "new", Content: "new content", Channel: domain.ChannelEmail,
	})

	require.NoError(t, err)
	assert.Equal(t, "new", result.Title)
	assert.Equal(t, "new", persistedTitle)
}

func TestNotificationService_Delete_OwnershipMismatchReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)

	ownerID := uuid.New()
	otherUserID := uuid.New()
	notificationID := uuid.New()

	repo.EXPECT().FindByID(gomock.Any(), notificationID).
		Return(&domain.Notification{ID: notificationID, UserID: ownerID}, nil)
	// repo.Delete is intentionally never EXPECTed: gomock fails the test if it's called.

	svc := service.NewNotificationService(repo, mocks.NewMockUserRepository(ctrl), service.SenderRegistry{})

	err := svc.Delete(context.Background(), otherUserID, notificationID)

	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestNotificationService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)

	userID := uuid.New()
	notificationID := uuid.New()

	repo.EXPECT().FindByID(gomock.Any(), notificationID).
		Return(&domain.Notification{ID: notificationID, UserID: userID}, nil)
	repo.EXPECT().Delete(gomock.Any(), notificationID).Return(nil)

	svc := service.NewNotificationService(repo, mocks.NewMockUserRepository(ctrl), service.SenderRegistry{})

	err := svc.Delete(context.Background(), userID, notificationID)
	require.NoError(t, err)
}

func TestNotificationService_ListByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockNotificationRepository(ctrl)

	userID := uuid.New()
	expected := []domain.Notification{{ID: uuid.New(), UserID: userID}}

	repo.EXPECT().FindByUserId(gomock.Any(), userID).Return(expected, nil)

	svc := service.NewNotificationService(repo, mocks.NewMockUserRepository(ctrl), service.SenderRegistry{})

	result, err := svc.ListByUser(context.Background(), userID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, expected[0].ID, result[0].ID)
}
