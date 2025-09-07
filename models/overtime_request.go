// internal/models/overtime_request.go

package models

import (
	"time"

	"github.com/uptrace/bun"
)

type OvertimeRequest struct {
	bun.BaseModel `bun:"table:overtime_requests,alias:or"`

	ID           int64     `bun:"id,pk,autoincrement"`
	Status       string    `bun:"status,notnull,default:'pending'"` // 'pending', 'approved', 'rejected'
	RequestTime  time.Time `bun:"request_time,notnull"`


	UserID int64 `bun:"user_id,notnull"`
	SlotID int64 `bun:"slot_id,notnull"`


	ReviewedBy *int64 `bun:"reviewed_by"`


	User *User         `bun:"rel:belongs-to,join:user_id=id"`
	Slot *OvertimeSlot `bun:"rel:belongs-to,join:slot_id=id"`
}