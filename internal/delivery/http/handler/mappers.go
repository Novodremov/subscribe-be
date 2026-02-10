package handler

import (
	"time"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/Novodremov/subscribe-be/internal/dto"

	"github.com/google/uuid"
)

func MapCreateDTOToCreateDomain(in dto.CreateSubscriptionRequest) (domain.CreateSubscription, error) {
	startDate, err := time.Parse("01-2006", in.StartDate)
	if err != nil {
		return domain.CreateSubscription{}, err
	}

	var endDate *time.Time
	if in.EndDate != "" {
		t, err := time.Parse("01-2006", in.EndDate)
		if err != nil {
			return domain.CreateSubscription{}, err
		}
		endDate = &t
	}

	userUUID, err := uuid.Parse(in.UserID)
	if err != nil {
		return domain.CreateSubscription{}, err
	}

	return domain.CreateSubscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      userUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func MapUpdateDTOToUpdateDomain(in dto.UpdateSubscriptionRequest) (domain.UpdateSubscription, error) {
	var startDate, endDate *time.Time

	if in.StartDate != nil {
		t, err := time.Parse("01-2006", *in.StartDate)
		if err != nil {
			return domain.UpdateSubscription{}, err
		}
		startDate = &t
	}

	if in.EndDate != nil && *in.EndDate != "" {
		t, err := time.Parse("01-2006", *in.EndDate)
		if err != nil {
			return domain.UpdateSubscription{}, err
		}
		endDate = &t
	}

	return domain.UpdateSubscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func MapDomainToCreateResponse(sub domain.Subscription) dto.CreateSubscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		s := sub.EndDate.Format("01-2006")
		endDate = &s
	}

	return dto.CreateSubscriptionResponse{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDate,
	}
}

func MapDomainToResponse(sub domain.Subscription) dto.SubscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		s := sub.EndDate.Format("01-2006")
		endDate = &s
	}

	return dto.SubscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDate,
		CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sub.UpdatedAt.Format(time.RFC3339),
	}
}

func MapDomainSubscriptionsToDTO(subs []domain.Subscription, totalCount int) dto.ListSubscriptionsResponse {
	res := make([]dto.SubscriptionResponse, 0, len(subs))

	for _, sub := range subs {
		res = append(res, MapDomainToResponse(sub))
	}

	return dto.ListSubscriptionsResponse{
		Subscriptions: res,
		TotalCount:    totalCount,
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
