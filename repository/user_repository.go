package repository

import (
	"context"
	"shiftdony/models"

	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *userRepository) GetUserByPersonnelCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("personnel_code = ?", code).
		Scan(ctx)
	return &user, err
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Relation("Team", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("manager_id")
		}).
		Where("u.id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
