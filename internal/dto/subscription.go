package dto

import "github.com/google/uuid"

// CreateSubscriptionRequest represents the payload to create a subscription
// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	// Name of the service to subscribe to
	ServiceName string `json:"service_name" example:"Netflix"`
	// Subscription price in rubles
	Price int `json:"price" example:"599"`
	// ID of the user creating the subscription
	UserID uuid.UUID `json:"user_id" example:"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"`
	// Subscription start date in MM-YYYY format
	StartDate string `json:"start_date" example:"01-2026"`
	// Optional subscription end date in MM-YYYY format
	EndDate string `json:"end_date" example:"12-2026"`
}

// UpdateSubscriptionRequest represents the payload to update a subscription
// swagger:model UpdateSubscriptionRequest
type UpdateSubscriptionRequest struct {
	// New service name (optional)
	ServiceName *string `json:"service_name,omitempty" example:"Spotify"`
	// New subscription price in rubles (optional)
	Price *int `json:"price,omitempty" example:"299"`
	// Updated start date in MM-YYYY format (optional)
	StartDate *string `json:"start_date,omitempty" example:"02-2026"`
	// Updated end date in MM-YYYY format (optional)
	EndDate *string `json:"end_date,omitempty" example:"11-2026"`
}

// SubscriptionResponse represents a subscription returned from the API
// swagger:model SubscriptionResponse
type SubscriptionResponse struct {
	// Unique subscription ID
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Name of the subscribed service
	ServiceName string `json:"service_name" example:"Netflix"`
	// Price in cents
	Price int `json:"price" example:"599"`
	// ID of the user owning the subscription
	UserID uuid.UUID `json:"user_id" example:"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"`
	// Subscription start month-year in MM-YYYY format
	StartDate string `json:"start_date" example:"01-2026"`
	// Optional subscription end month-year in MM-YYYY format
	EndDate *string `json:"end_date,omitempty" example:"12-2026"`
	// Timestamp when subscription was created
	CreatedAt string `json:"created_at" example:"2026-03-05T10:30:35Z"`
	// Timestamp when subscription was last updated
	UpdatedAt string `json:"updated_at" example:"2026-03-05T10:30:35Z"`
}

// ListSubscriptionsResponse represents a paginated list of subscriptions
// swagger:model ListSubscriptionsResponse
type ListSubscriptionsResponse struct {
	// List of subscriptions
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	// Total number of subscriptions
	TotalCount int `json:"total_count" example:"42"`
}

// SubscriptionsTotalCostResponse represents aggregated subscription cost
// swagger:model SubscriptionsTotalCostResponse
type SubscriptionsTotalCostResponse struct {
	// Total cost of subscriptions in rubles
	TotalCost int64 `json:"total_cost" example:"504328"`
}
