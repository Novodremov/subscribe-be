package repo

import (
	"context"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/google/uuid"
)

type ISubscriptionRepo interface {
	Create(ctx context.Context, in *domain.CreateSubscription) (*domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]domain.Subscription, int, error)
	ListFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]domain.Subscription, int, error)
}
