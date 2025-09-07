package models

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64  `bun:"id,pk,autoincrement"`
	PersonnelCode string `bun:"personnel_code,unique,notnull"`
	FullName      string `bun:"full_name,notnull"`
	PasswordHash  string `bun:"password_hash,notnull"`
	Role          string `bun:"role,notnull"`
	WorkHours     string `bun:"work_hours"`

	TeamID int64 `bun:"team_id,notnull"`
	Team   *Team `bun:"rel:belongs-to,join:team_id=id"`
}
