package repo

import (
	"context"
	"fmt"

	sqlc "github.com/Novodremov/subscribe-be/internal/db/sqlc_generated"
	"github.com/Novodremov/subscribe-be/internal/domain"

	"github.com/google/uuid"
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
	var pgEndDate pgtype.Date
	if in.EndDate != nil {
		pgEndDate = pgtype.Date{
			Time:  *in.EndDate,
			Valid: true,
		}
	} else {
		pgEndDate = pgtype.Date{Valid: false}
	}
	params := sqlc.CreateSubscriptionParams{
		ServiceName: in.ServiceName,
		Price:       int32(in.Price),
		UserID:      uuidToPgtype(in.UserID),
		StartDate:   timeToPGDate(in.StartDate),
		EndDate:     pgEndDate,
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
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return mapSQLCToDomain(sub)
}

func (r *subscriptionRepo) Update(ctx context.Context, in *domain.UpdateSubscription) (*domain.Subscription, error) {
	params := sqlc.UpdateSubscriptionParams{
		ID:          uuidToPgtype(in.ID),
		ServiceName: derefString(in.ServiceName),
		Price:       int32(derefInt(in.Price)),
	}

	if in.StartDate != nil {
		params.StartDate = timeToPGDate(*in.StartDate)
	}
	if in.EndDate != nil {
		params.EndDate = timeToPGDate(*in.EndDate)
	}

	sub, err := r.q.UpdateSubscription(ctx, r.conn, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return mapSQLCToDomain(sub)
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	pgID := uuidToPgtype(id)
	if err := r.q.DeleteSubscription(ctx, r.conn, pgID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepo) List(ctx context.Context, limit, offset int) ([]domain.Subscription, int, error) {
	params := sqlc.ListSubscriptionsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	subs, err := r.q.ListSubscriptions(ctx, r.conn, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	domainSubs := make([]domain.Subscription, 0, len(subs))
	for _, s := range subs {
		ds, err := mapSQLCToDomain(s)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to map subscription: %w", err)
		}
		domainSubs = append(domainSubs, *ds)
	}

	return domainSubs, len(domainSubs), nil
}

func (r *subscriptionRepo) ListFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]domain.Subscription, int, error) {
	var pgID pgtype.UUID
	if userID != nil {
		pgID = uuidToPgtype(*userID)
	} else {
		pgID = pgtype.UUID{Valid: false}
	}

	svc := ""
	if serviceName != nil {
		svc = *serviceName
	}

	params := sqlc.ListSubscriptionsFilteredParams{
		Column1: pgID,
		Column2: svc,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	subs, err := r.q.ListSubscriptionsFiltered(ctx, r.conn, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list filtered subscriptions: %w", err)
	}

	domainSubs := make([]domain.Subscription, 0, len(subs))
	for _, s := range subs {
		ds, err := mapSQLCToDomain(s)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to map subscription: %w", err)
		}
		domainSubs = append(domainSubs, *ds)
	}

	return domainSubs, len(domainSubs), nil
}
