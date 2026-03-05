package service

import (
	"context"
	"time"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/google/uuid"
)

type ISubscriptionService interface {
	CreateSubscription(ctx context.Context, in *domain.CreateSubscription) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, limit, offset int) ([]domain.Subscription, int, error)
	SubscriptionsTotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int64, error)
}
