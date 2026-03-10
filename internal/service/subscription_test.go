package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/Novodremov/subscribe-be/internal/service"
	mock_repo "github.com/Novodremov/subscribe-be/internal/repo/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionService_CreateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	in := &domain.CreateSubscription{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
		EndDate:     nil,
	}

	expected := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	}

	mockRepo.EXPECT().
		Create(ctx, in).
		Return(expected, nil).
		Times(1)

	actual, err := svc.CreateSubscription(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscriptionService_CreateSubscription_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	in := &domain.CreateSubscription{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
		EndDate:     nil,
	}

	repoErr := errors.New("repo fail")
	mockRepo.EXPECT().
		Create(ctx, in).
		Return(nil, repoErr).
		Times(1)

	actual, err := svc.CreateSubscription(ctx, in)
	require.Nil(t, actual)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repo create subscription failed")
	assert.ErrorIs(t, err, repoErr)
}

func TestSubscriptionService_GetSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	expected := &domain.Subscription{
		ID:          id,
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.EXPECT().
		Get(ctx, id).
		Return(expected, nil).
		Times(1)

	actual, err := svc.GetSubscription(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscriptionService_GetSubscription_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	repoErr := errors.New("repo fail")
	mockRepo.EXPECT().
		Get(ctx, id).
		Return(nil, repoErr).
		Times(1)

	actual, err := svc.GetSubscription(ctx, id)
	require.Nil(t, actual)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repo get subscription failed")
	assert.ErrorIs(t, err, repoErr)
}

func TestSubscriptionService_UpdateSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	update := &domain.UpdateSubscription{
		ID:          uuid.New(),
		ServiceName: ptrString("Netflix"),
		Price:       ptrInt(699),
	}

	expected := &domain.Subscription{
		ID:          update.ID,
		ServiceName: *update.ServiceName,
		Price:       *update.Price,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.EXPECT().
		Update(ctx, update).
		Return(expected, nil).
		Times(1)

	actual, err := svc.UpdateSubscription(ctx, update)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscriptionService_UpdateSubscription_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	update := &domain.UpdateSubscription{
		ID:          uuid.New(),
		ServiceName: ptrString("Netflix"),
	}

	repoErr := errors.New("repo update fail")
	mockRepo.EXPECT().
		Update(ctx, update).
		Return(nil, repoErr).
		Times(1)

	actual, err := svc.UpdateSubscription(ctx, update)
	require.Nil(t, actual)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repo update subscription failed")
	assert.ErrorIs(t, err, repoErr)
}

func TestSubscriptionService_DeleteSubscription_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(nil).
		Times(1)

	err := svc.DeleteSubscription(ctx, id)
	require.NoError(t, err)
}

func TestSubscriptionService_DeleteSubscription_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	id := uuid.New()
	repoErr := errors.New("repo delete fail")

	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(repoErr).
		Times(1)

	err := svc.DeleteSubscription(ctx, id)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repo delete subscription failed")
	assert.ErrorIs(t, err, repoErr)
}

func TestSubscriptionService_ListSubscriptions_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subs := []domain.Subscription{
		{ID: uuid.New(), ServiceName: "Netflix", Price: 599},
	}
	total := 1

	mockRepo.EXPECT().
		List(ctx, 10, 0).
		Return(subs, nil).
		Times(1)
	mockRepo.EXPECT().
		TotalCount(ctx).
		Return(total, nil).
		Times(1)

	gotSubs, gotTotal, err := svc.ListSubscriptions(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, subs, gotSubs)
	assert.Equal(t, total, gotTotal)
}

func TestSubscriptionService_ListSubscriptions_RepoListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	repoErr := errors.New("repo list fail")

	mockRepo.EXPECT().
		List(ctx, 10, 0).
		Return(nil, repoErr).
		Times(1)

	gotSubs, gotTotal, err := svc.ListSubscriptions(ctx, 10, 0)
	require.Error(t, err)
	assert.Nil(t, gotSubs)
	assert.Equal(t, 0, gotTotal)
	assert.Contains(t, err.Error(), "repo list subscriptions failed")
}

func TestSubscriptionService_ListSubscriptions_RepoTotalCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subs := []domain.Subscription{
		{ID: uuid.New(), ServiceName: "Netflix", Price: 599},
	}
	totalErr := errors.New("repo total count fail")

	mockRepo.EXPECT().
		List(ctx, 10, 0).
		Return(subs, nil).
		Times(1)
	mockRepo.EXPECT().
		TotalCount(ctx).
		Return(0, totalErr).
		Times(1)

	gotSubs, gotTotal, err := svc.ListSubscriptions(ctx, 10, 0)
	require.Error(t, err)
	assert.Nil(t, gotSubs)
	assert.Equal(t, 0, gotTotal)
	assert.Contains(t, err.Error(), "failed to get total count of subscriptions")
}

func TestSubscriptionService_SubscriptionsTotalCost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	serviceName := "Netflix"
	startDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	expectedTotal := int64(1234)

	mockRepo.EXPECT().
		TotalCost(ctx, &userID, &serviceName, &startDate, &endDate).
		Return(expectedTotal, nil).
		Times(1)

	total, err := svc.SubscriptionsTotalCost(ctx, &userID, &serviceName, &startDate, &endDate)
	require.NoError(t, err)
	assert.Equal(t, expectedTotal, total)
}

func TestSubscriptionService_SubscriptionsTotalCost_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockISubscriptionRepo(ctrl)
	svc := service.NewSubscriptionService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	repoErr := errors.New("repo total cost fail")

	mockRepo.EXPECT().
		TotalCost(ctx, &userID, nil, nil, nil).
		Return(int64(0), repoErr).
		Times(1)

	total, err := svc.SubscriptionsTotalCost(ctx, &userID, nil, nil, nil)
	require.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Contains(t, err.Error(), "repo total cost calculation failed")
}

// хелперы для указателей
func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }