package service

import (
	"context"
	"database/sql"
	"shiftdony/config"
	postgres "shiftdony/database"
	"shiftdony/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo postgres.UserRepository
}

func NewUserService(userRepo postgres.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(ctx context.Context, personnelCode, fullName, password string, teamID int64) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternalServer
	}
	newUser := models.User{
		PersonnelCode: personnelCode,
		FullName: fullName,
		PasswordHash: string(hashedPassword),
		Role: "user",
		TeamID: teamID,
		WorkHours: "9-17",
	}
	if err := s.userRepo.CreateUser(ctx, &newUser); err != nil {
		if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() {
			return ErrPersonnelCodeExists
		}
		return ErrInternalServer
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, personnelCode, password string) (string, error) {
    user, err := s.userRepo.GetUserByPersonnelCode(ctx, personnelCode)
    if err != nil {
        return "", ErrInvalidCredentials 
    }

   
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return "", ErrInvalidCredentials
    }

    claims := jwt.MapClaims{
        "sub":  user.ID,
        "role": user.Role,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
        "iat":  time.Now().Unix(),
    }

   
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
    tokenString, err := token.SignedString([]byte(config.C.JWT.Secret))
    if err != nil {
        return "", ErrInternalServer
    }

    return tokenString, nil
}


func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.User, error){
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}
	user.PasswordHash = ""
	return user, nil
}