package routes

import (
	"shiftdony/handlers"
	"shiftdony/middleware"
	"shiftdony/repository"
	"shiftdony/service"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func SetupRouter(db *bun.DB) *gin.Engine {
	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	overtimeRepo := repository.NewOvertimeRepository(db)

	userService := service.NewUserService(userRepo)
	overtimeService := service.NewOvertimeService(overtimeRepo)

	userHandler := handlers.NewUserHandler(userService)
	overtimeHandler := handlers.NewOvertimeHandler(overtimeService)
	
	//Public Routes
	// Public Routes
	api := router.Group("/api")
	{
		api.POST("/register", userHandler.RegisterUser)
		api.POST("/login", userHandler.Login)
	}

	// Protected Routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/overtime/available", overtimeHandler.GetAvailableOvertimeSlots)
		protected.POST("/requests/append", overtimeHandler.CreateOvertimeRequest)
		protected.GET("/my-requests", overtimeHandler.GetMyOvertimeRequests)

		// Admins Routes
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.AdminMiddleware())
		{
			adminRoutes.POST("/overtime", overtimeHandler.CreateOvertimeSlot)
			adminRoutes.GET("/overtime", overtimeHandler.GetOvertimeSlots)
			adminRoutes.GET("/requests", overtimeHandler.GetAllOvertimeRequests)
			adminRoutes.PATCH("/requests/:id", overtimeHandler.UpdateOvertimeReqStatus)
			adminRoutes.GET("/reports/csv", overtimeHandler.ExportApprovedRequestsAsCSV)
		}
	}

	return router
}
