package routes

import (
	postgres "shiftdony/database"
	"shiftdony/handlers"
	"shiftdony/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *postgres.Postgres) *gin.Engine {
	router := gin.Default()

	userHandler := handlers.NewUserHandler(db)
	overtimeHandler := handlers.NewOvertimeHandler(db)

	//Public Routes
	api := router.Group("/api")
	{
		api.POST("/register", userHandler.RegisterUser)
		api.POST("/login", userHandler.Login)
	}

	//Protected Routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/overtime/available", overtimeHandler.GetAvailableOvertimeSlots)
		protected.POST("/requests", overtimeHandler.CreateOvertimeRequest)
		protected.GET("/my-requests", overtimeHandler.GetMyOvertimeRequests)

		//Admins Routes
		adminRoutrs := protected.Group("/admin")
		adminRoutrs.Use(middleware.AdminMiddleware())
		{
			adminRoutrs.POST("/overtime", overtimeHandler.CreateOvertimeSlot)
			adminRoutrs.GET("/overtime", overtimeHandler.GetOvertimeSlots)
			adminRoutrs.GET("/requests", overtimeHandler.GetAllOvertimeRequests)
			adminRoutrs.PATCH("/requests/:id", overtimeHandler.UpdateOvertimeReqStatus)
			adminRoutrs.GET("/reports/csv", overtimeHandler.ExportApprovedRequestsAsCSV)
		}
	}

	return router
}
