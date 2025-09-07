package models

import "github.com/uptrace/bun"

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Name      string `bun:"name,notnull"`
	ManagerID int64  `bun:"manager_id"`
}
