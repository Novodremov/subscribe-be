package dto

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Netflix"`
	Price       int    `json:"price" example:"599"`
	UserID      string `json:"user_id" example:"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"`
	StartDate   string `json:"start_date" example:"01-2026"`
	EndDate     string `json:"end_date" example:"12-2026"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name" example:"Spotify"`
	Price       *int    `json:"price" example:"299"`
	StartDate   *string `json:"start_date" example:"02-2026"`
	EndDate     *string `json:"end_date" example:"11-2026"`
}


type ListSubscriptionsRequest struct {
	UserID      *string `json:"user_id"`
	ServiceName *string `json:"service_name"`
	Limit       int     `json:"limit"`
	Offset      int     `json:"offset"`
}

type CreateSubscriptionResponse struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	TotalCount    int                    `json:"total_count"`
}
