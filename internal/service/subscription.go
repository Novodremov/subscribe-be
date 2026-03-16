package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/Novodremov/subscribe-be/internal/repo"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repo.ISubscriptionRepo
}

// NewSubscriptionService создает новый экземпляр SubscriptionService с указанным репозиторием. Используется для бизнес-логики работы с подписками.
func NewSubscriptionService(r repo.ISubscriptionRepo) ISubscriptionService {
	return &SubscriptionService{
		repo: r,
	}
}

// CreateSubscription создает новую подписку через репозиторий. Возвращает созданную подписку или ошибку при сохранении.
func (s *SubscriptionService) CreateSubscription(ctx context.Context, in *domain.CreateSubscription) (*domain.Subscription, error) {
	// Здесь могла бы быть ваша бизнес-логика
	sub, err := s.repo.Create(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("repo create subscription failed: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return sub, nil
}

// GetSubscription возвращает подписку по ID или ошибку, если не найдена.
func (s *SubscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	// Здесь могла бы быть ваша бизнес-логика
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo get subscription failed: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return sub, nil
}

// UpdateSubscription обновляет существующую подписку и возвращает её после изменений.
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error) {
	// Здесь могла бы быть ваша бизнес-логика
	sub, err := s.repo.Update(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("repo update subscription failed: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return sub, nil
}

// DeleteSubscription удаляет подписку по ID.
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo delete subscription failed: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return nil
}

// ListSubscriptions возвращает список подписок с пагинацией и их общее количество.
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, limit, offset int) ([]domain.Subscription, int, error) {
	// Здесь могла бы быть ваша бизнес-логика
	subs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list subscriptions failed: %w", err)
	}
	total, err := s.repo.TotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count of subscriptions: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return subs, total, nil
}

// SubscriptionsTotalCost вычисляет суммарную стоимость подписок с учетом фильтров.
func (s *SubscriptionService) SubscriptionsTotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int64, error) {
	// Здесь могла бы быть ваша бизнес-логика
	total, err := s.repo.TotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("repo total cost calculation failed: %w", err)
	}
	// Здесь могла бы быть ваша бизнес-логика
	return total, nil
}
