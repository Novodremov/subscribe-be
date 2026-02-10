package handler

import "errors"

var (
	ErrValidation       = errors.New("validation")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrEmptyUserID      = errors.New("user_id is required")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
