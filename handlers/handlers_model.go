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


type SlotResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Capacity  int64     `json:"capacity"`
	Status    string    `json:"status"`
	Creator   string    `json:"creator"` 
}




