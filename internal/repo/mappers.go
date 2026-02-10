package repo

import (
	"fmt"
	"time"

	sqlc "github.com/Novodremov/subscribe-be/internal/db/sqlc_generated"
	"github.com/Novodremov/subscribe-be/internal/domain"
)

func mapSQLCToDomain(s *sqlc.Subscription) (*domain.Subscription, error) {
	id, err := pgtypeToUUID(s.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid subscription ID: %w", err)
	}

	userID, err := pgtypeToUUID(s.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var endDate *time.Time
	if s.EndDate.Valid {
		t := time.Time(s.EndDate.Time)
		endDate = &t
	}

	createdAt, err := pgTimestamptzToTime(s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at: %w", err)
	}
	
	updatedAt, err := pgTimestamptzToTime(s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at: %w", err)
	}
	return &domain.Subscription{
		ID:          id,
		ServiceName: s.ServiceName,
		Price:       int(s.Price),
		UserID:      userID,
		StartDate:   s.StartDate.Time,
		EndDate:     endDate,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
