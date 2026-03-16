package handler

import (
	"time"

	"github.com/Novodremov/subscribe-be/internal/domain"
	"github.com/Novodremov/subscribe-be/internal/dto"
)

// MapCreateDTOToCreateDomain преобразует dto.CreateSubscriptionRequest в domain.CreateSubscription.
func MapCreateDTOToCreateDomain(in dto.CreateSubscriptionRequest) (domain.CreateSubscription, error) {
	startDate, err := time.Parse(DateShortLayout, in.StartDate)
	if err != nil {
		return domain.CreateSubscription{}, err
	}

	var endDate *time.Time
	if in.EndDate != "" {
		t, err := time.Parse(DateShortLayout, in.EndDate)
		if err != nil {
			return domain.CreateSubscription{}, err
		}
		endDate = &t
	}

	return domain.CreateSubscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      in.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

// MapCreateDTOToCreateDomain преобразует dto.UpdateSubscriptionRequest в domain.UpdateSubscription.
func MapUpdateDTOToUpdateDomain(in dto.UpdateSubscriptionRequest) (domain.UpdateSubscription, error) {
	var startDate, endDate *time.Time

	if in.StartDate != nil {
		t, err := time.Parse(DateShortLayout, *in.StartDate)
		if err != nil {
			return domain.UpdateSubscription{}, err
		}
		startDate = &t
	}

	if in.EndDate != nil && *in.EndDate != "" {
		t, err := time.Parse(DateShortLayout, *in.EndDate)
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

// MapDomainToResponse преобразует domain.Subscription в dto.SubscriptionResponse.
func MapDomainToResponse(sub domain.Subscription) dto.SubscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		s := sub.EndDate.Format(DateShortLayout)
		endDate = &s
	}

	return dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format(DateShortLayout),
		EndDate:     endDate,
		CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sub.UpdatedAt.Format(time.RFC3339),
	}
}

// MapDomainToResponse преобразует []domain.Subscription в dto.ListSubscriptionsResponse.
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
