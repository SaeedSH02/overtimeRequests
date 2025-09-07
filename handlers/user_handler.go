package handlers

import (
	"context"
	"net/http"
	"time"

	"shiftdony/config"
	postgres "shiftdony/database"
	log "shiftdony/logs"
	"shiftdony/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *postgres.Postgres
}

func NewUserHandler(db *postgres.Postgres) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid Input",
			"reason":  "INVALID_INPUT"})
		log.Gl.Error("Failed to register:", zap.Error(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Error("Failed to hash password", zap.Error(err))
		return
	}

	newUser := models.User{
		PersonnelCode: input.PersonnelCode,
		FullName:      input.FullName,
		PasswordHash:  string(hashedPassword),
		Role:          "user",
		TeamID:        input.TeamID,
		WorkHours:     "9-17",
	}

	_, err = h.db.DB().NewInsert().Model(&newUser).Exec(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to creat user",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("user can't Registred, inputs: ",
			zap.String("PersonnelCode", newUser.PersonnelCode),
			zap.String("Passwords", newUser.PasswordHash),
			zap.String("userName", newUser.FullName))
		return
	}
	log.Gl.Info("user Registred to pannel: ",
		zap.Int64("userID", newUser.ID),
		zap.String("PersonnelCode", newUser.PersonnelCode),
		zap.String("userName", newUser.FullName))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	err := h.db.DB().
		NewSelect().
		Model(&user).
		Where("personnel_code = ?", input.PersonnelCode).
		Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid personnel code or password",
			"reason":  "Unauthorized"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid personnel code or password",
			"reason":  "Unauthorized"})
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expiration time
		"iat":  time.Now().Unix(),
	}

	//Create jwt with HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.C.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "SERVER_ERROR"})
		log.Gl.Info("error while making JWT" + err.Error())
		return
	}

	log.Gl.Info("user loged to pannel: ", zap.Int64("userID", user.ID),
		zap.String("userRole", user.Role), zap.String("userName", user.FullName))

	c.JSON(http.StatusOK, gin.H{"success": true, "token": tokenString})
}

// Get user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Internal Server Error",
			"reason":  "Unauthorized"})
		return
	}

	var user models.User
	err := h.db.DB().NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Not Found",
			"reason":  "NOT_FOUND"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"success": true, "profile": user})

}
