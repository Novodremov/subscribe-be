package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/Novodremov/subscribe-be/internal/domain"
)

type ISubscriptionService interface {
	CreateSubscription(ctx context.Context, in *domain.CreateSubscription) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, limit, offset int) ([]domain.Subscription, int, error)
	ListSubscriptionsFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]domain.Subscription, int, error)
}
