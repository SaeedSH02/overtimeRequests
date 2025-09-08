package handlers

import (
	"net/http"

	log "shiftdony/logs"
	"shiftdony/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input RegisterInput

	//Read inputs
	if err := c.ShouldBindJSON(&input); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid input data", "INVALID_INPUT")
		return
	}

	//Service Call
	err := h.userService.Register(
		c.Request.Context(),
		input.PersonnelCode,
		input.FullName,
		input.Password,
		input.TeamID,
	)
	if err != nil {
		if err == service.ErrPersonnelCodeExists {
			SendErrorResponse(c, http.StatusConflict, "A user with this personnel code already exists", "ALREADY_EXISTS")
			log.Gl.Info("A user with this personnel code already exists", zap.String("Personnel Code", input.PersonnelCode))

		} else {
			SendErrorResponse(c, http.StatusInternalServerError, "Failed to register user", "SERVER_ERROR")
			log.Gl.Error("Failed to register user", zap.Error(err))
		}
		return
	}

	//Success
	SendSuccessResponse(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})

}

func (h *UserHandler) Login(c *gin.Context) {
	var input LoginInput
	//Read inputs
	if err := c.ShouldBindJSON(&input); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "Invalid input data", "INVALID_INPUT")
		log.Gl.Info("Invalid input data for login", zap.String("Personnel Code", input.PersonnelCode))
		return
	}

	//Service Call
	token, err := h.userService.Login(c.Request.Context(), input.PersonnelCode, input.Password)

	if err != nil {
		if err == service.ErrInvalidCredentials {
			SendErrorResponse(c, http.StatusUnauthorized, "Invalid personnel code or password", "INVALID_CREDENTIALS")
			log.Gl.Info("Invalid personnel code or password", zap.String("Personnel Code", input.PersonnelCode))
		} else {
			SendErrorResponse(c, http.StatusInternalServerError, "Could not process login", "SERVER_ERROR")
			log.Gl.Error("Could not process login", zap.String("Personnel Code", input.PersonnelCode), zap.Error(err))
		}
		return
	}

	//Success
	SendSuccessResponse(c, http.StatusOK, gin.H{
		"token": token,
	})
}

// Get user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		log.Gl.Info("User not authenticated")
		return
	}

	userID, ok := userIDVal.(float64)
	if !ok {
		SendErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format in token", "SERVER_ERROR")
		return
	}

	userProfile, err := h.userService.GetProfile(c.Request.Context(), int64(userID))
	if err != nil {
		if err == service.ErrUserNotFound {
			SendErrorResponse(c, http.StatusNotFound, "User profile not found", "NOT_FOUND")
			log.Gl.Info("User profile not found", zap.Error(err))
		} else {
			SendErrorResponse(c, http.StatusInternalServerError, "Could not retrieve profile", "SERVER_ERROR")
			log.Gl.Error("Could not retrieve profile", zap.Error(err))
		}
		return

	}

	SendSuccessResponse(c, http.StatusOK, userProfile)
}
