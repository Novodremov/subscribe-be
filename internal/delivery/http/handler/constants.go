package handler

const (
	DefaultOffset = 0
	DefaultLimit  = 50
	MinLimit      = 1
	MaxLimit      = 100
)

const (
	DateShortLayout = "01-2006"
	DateLayout      = "02-01-2006"
)

const (
	ErrMsgInvalidBody         = "invalid body"
	ErrMsgNoFieldsToUpdate    = "no fields to update"
	ErrMsgEmptyID             = "subscription_id is required"
	ErrMsgInvalidData         = "invalid data"
	ErrMsgBadRequest          = "bad request"
	ErrMsgInternalServerError = "internal server error"
)
