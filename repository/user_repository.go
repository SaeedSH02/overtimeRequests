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
	_, err := r.db.NewSelect().Model(user).Exec(ctx)
	return err
}

func (r *userRepository) GetUserByPersonnelCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("perssonel_code = ?", code).
		Scan(ctx)
	return &user, err
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Scan(ctx)
	return &user, err
}

