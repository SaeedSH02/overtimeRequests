package handlers

import "time"

//handler input models

type RegisterInput struct {
	PersonnelCode string `json:"personnel_code" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	TeamID        int64  `json:"team_id" binding:"required"`
}

type LoginInput struct {
	PersonnelCode string `json:"personnel_code" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type CreateOvertimeInput struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Capacity  int       `json:"capacity" binding:"required"`
}

type CreateRequestInput struct {
	SlotID int64 `json:"slot_id" binding:"required"`
}

type UpdateRequestStatusInput struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}




// SERVER_ERROR
// CAPTCHA_ERROR
// INVALID_CREDENTIALS
// USER_ALREADY_LOGGED_IN
// TOKEN_GENERATION_FAILED
// NOT_WHITELISTED
// USER_ALREADY_EXISTS
// SIGNED_UP
// TOKEN_REVOKED
// INVALID_TOKEN
// TOKEN_UNVERIFIED
// NOT_AUTHORIZED
// INVALID_TOKEN
// USER_NOT_FOUND

