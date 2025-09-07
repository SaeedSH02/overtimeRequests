package postgres

import (
	"context"
	"shiftdony/models"
)

type DB interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)

	
	DeleteChat(ctx context.Context, chatCtxID int64) error
	// AddTicket(ctx context.Context, ticket *models.Ticket) error
	// UpdateTicketMessage(ctx context.Context, userID int64, message string) (bool, error)
}
