package app

import "github.com/gin-gonic/gin"

// RegisterRoutes registers a minimal set of routes for the application.
func RegisterRoutes(r *gin.Engine, a *App) {
	// Mount socket server into Gin
	r.GET("/socket.io/*any", gin.WrapH(a.SocketServer))
	r.POST("/socket.io/*any", gin.WrapH(a.SocketServer))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "FMIIS backend server is running"})
	})

	api := r.Group("/api/fmiis-backend/v001")
	{
		// PUBLIC routes
		api.POST("/login", a.UserHandler.Login)
		// api.POST("/register", a.UserHandler.Register)

		// PROTECTED routes
		protectedApi := api.Group("")
		protectedApi.Use(a.AuthMiddleware.Handle())
		{
			protectedApi.GET("/auth/me")
			protectedApi.POST("/register", a.AuthMiddleware.RequireRole("gm", "md"), a.UserHandler.Register)
			protectedApi.POST("/send-notification", a.NotificationHandler.SendNotification)
			protectedApi.POST("/logout", a.UserHandler.Logout)
			protectedApi.POST("/gm-update", a.AuthMiddleware.RequireRole("gm", "md"), a.UserHandler.GmUpdate)
			protectedApi.POST("/normal-update", a.UserHandler.NormalUpdate)
			protectedApi.POST("/delete/:id", a.AuthMiddleware.RequireRole("gm", "md"), a.UserHandler.GmDelete)
			protectedApi.POST("/add-attendance-record", a.AuthMiddleware.RequireRole("gm", "md"), a.AttendanceHandler.CreateAttendance)
			protectedApi.POST("/request-leave", a.LeaveHandler.CreateLeaveRequest)
			protectedApi.POST("/leave-request-approval", a.AuthMiddleware.RequireRole("gm"), a.LeaveHandler.LeaveRequestGmApproval)
			protectedApi.GET("/get-all-employees", a.AuthMiddleware.RequireRole("gm", "md", "hr"), a.UserHandler.GetAllEmployees)
			protectedApi.GET("/get-single-employee/:id", a.AuthMiddleware.RequireRole("gm", "md", "hr"), a.UserHandler.GetSingleEmployee)
			protectedApi.GET("/get-attendance-data")
			// Other protected routes
		}
	}
}
