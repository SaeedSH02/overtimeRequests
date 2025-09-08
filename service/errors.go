package service

import "errors"

// Custom errors for business logic
var (
	ErrInvalidCredentials  = errors.New("invalid personnel code or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrPersonnelCodeExists = errors.New("personnel code already exists")

	ErrSlotNotFound   = errors.New("slot not found or is not open for requests")
	ErrAlreadyApplied = errors.New("you have already applied for this slot")
	ErrSlotIsFull     = errors.New("this overtime slot is already full")
	ErrRequestNotFound = errors.New("request not found")

	ErrInternalServer = errors.New("internal server error")
)