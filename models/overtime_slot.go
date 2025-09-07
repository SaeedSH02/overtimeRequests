// internal/models/overtime_slot.go

package models

import (
	"time"

	"github.com/uptrace/bun"
)

type OvertimeSlot struct {
	bun.BaseModel `bun:"table:overtime_slots,alias:os"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Title     string    `bun:"title,notnull"`
	StartTime time.Time `bun:"start_time,notnull"`
	EndTime   time.Time `bun:"end_time,notnull"`
	Capacity  int64     `bun:"capacity,notnull"`
	Status    string    `bun:"status,notnull,default:'open'"`

	CreatedBy int64 `bun:"created_by,notnull"`
	Creator   *User `bun:"rel:belongs-to,join:created_by=id"`
}
