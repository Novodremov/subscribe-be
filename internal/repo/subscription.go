package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	sqlc "github.com/Novodremov/subscribe-be/internal/db/sqlc_generated"
	"github.com/Novodremov/subscribe-be/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepo struct {
	q    sqlc.Querier
	conn sqlc.DBTX
}

func NewSubscriptionRepo(pool *pgxpool.Pool) ISubscriptionRepo {
	return &subscriptionRepo{
		q:    sqlc.New(),
		conn: pool,
	}
}

func (r *subscriptionRepo) Create(ctx context.Context, in *domain.CreateSubscription) (*domain.Subscription, error) {
	params := sqlc.CreateSubscriptionParams{
		ServiceName: in.ServiceName,
		Price:       int32(in.Price),
		UserID:      uuidToPgtype(in.UserID),
		StartDate:   timePtrToPGDate(&in.StartDate),
		EndDate:     timePtrToPGDate(in.EndDate),
	}

	sub, err := r.q.CreateSubscription(ctx, r.conn, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return mapSQLCToDomain(sub)
}

func (r *subscriptionRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	sub, err := r.q.GetSubscription(ctx, r.conn, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return mapSQLCToDomain(sub)
}

func (r *subscriptionRepo) Update(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error) {
	params := sqlc.UpdateSubscriptionParams{
		ID:          uuidToPgtype(in.ID),
		ServiceName: derefString(in.ServiceName),
		Price:       int32(derefInt(in.Price)),
		StartDate:   timePtrToPGDate(in.StartDate),
		EndDate:     timePtrToPGDate(in.EndDate),
	}

	sub, err := r.q.UpdateSubscription(ctx, r.conn, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return mapSQLCToDomain(sub)
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	pgID := uuidToPgtype(id)

	rows, err := r.q.DeleteSubscription(ctx, r.conn, pgID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	if rows == 0 {
			return ErrSubscriptionNotFound
		}
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, limit, offset int) ([]domain.Subscription, error) {
	params := sqlc.ListSubscriptionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	subs, err := r.q.ListSubscriptions(ctx, r.conn, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	domainSubs := make([]domain.Subscription, 0, len(subs))
	for _, s := range subs {
		ds, err := mapSQLCToDomain(s)
		if err != nil {
			return nil, fmt.Errorf("failed to map subscription: %w", err)
		}
		domainSubs = append(domainSubs, *ds)
	}

	return domainSubs, nil
}

func (r *subscriptionRepo) TotalCount(ctx context.Context) (int, error) {
	total, err := r.q.TotalSubscriptions(ctx, r.conn)
	if err != nil {
		return 0, fmt.Errorf("failed to get total number of subscriptions: %w", err)
	}
	return int(total), nil
}

func (r *subscriptionRepo) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate *time.Time) (int64, error) {
	var pgID pgtype.UUID
	if userID != nil {
		pgID = uuidToPgtype(*userID)
	} else {
		pgID = pgtype.UUID{Valid: false}
	}

	params := sqlc.SubscriptionsTotalCostParams{
		UserID:      pgID,
		ServiceName: serviceName,
		StartDate:   timePtrToPGDate(startDate),
		EndDate:     timePtrToPGDate(endDate),
	}

	total, err := r.q.SubscriptionsTotalCost(ctx, r.conn, params)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost of subscriptions: %w", err)
	}

	return total, nil
}
