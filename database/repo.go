package postgres

import (
	"context"
	"shiftdony/models"
)

// UserRepository defines the methods for interacting with user data.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByPersonnelCode(ctx context.Context, code string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
}

// OvertimeRepository defines the methods for interacting with overtime data.

type OvertimeRepository interface {
	CreateOvertimeSlot(ctx context.Context, slot *models.OvertimeSlot) error
	GetOvertimeSlots(ctx context.Context) ([]models.OvertimeSlot, error)
	GetAvailableOvertimeSlots(ctx context.Context) ([]models.OvertimeSlot, error)
	GetOvertimeSlotByID(ctx context.Context, slotID int64) (*models.OvertimeSlot, error)
	UserHasPendingRequestForSlot(ctx context.Context, userID, slotID int64) (bool, error)
	CountApprovedRequestsForSlot(ctx context.Context, slotID int64) (int, error)
	CreateOvertimeRequest(ctx context.Context, req *models.OvertimeRequest) error
	UpdateOvertimeSlot(ctx context.Context, slot *models.OvertimeSlot) error
	GetMyOvertimeRequests(ctx context.Context, userID int64) ([]models.OvertimeRequest, error)
	GetAllOvertimeRequests(ctx context.Context) ([]models.OvertimeRequest, error)
	GetOvertimeRequestByID(ctx context.Context, requestID int64) (*models.OvertimeRequest, error)
	UpdateOvertimeRequest(ctx context.Context, req *models.OvertimeRequest) error
	GetApprovedRequests(ctx context.Context) ([]models.OvertimeRequest, error)
}
