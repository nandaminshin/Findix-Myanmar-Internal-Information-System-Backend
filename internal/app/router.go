package app

import "github.com/gin-gonic/gin"

// RegisterRoutes registers a minimal set of routes for the application.
func RegisterRoutes(r *gin.Engine, a *App) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "FMIIS backend server is running"})
	})

	api := r.Group("/api/v1")
	{
		// PUBLIC routes
		api.POST("/login", a.UserHandler.Login)
		// api.POST("/register", a.UserHandler.Register)

		// PROTECTED routes
		protectedApi := api.Group("")
		protectedApi.Use(a.AuthMiddleware.Handle())
		{
			protectedApi.GET("/auth/me")
			protectedApi.POST("/register", a.AuthMiddleware.RequireRole("gm"), a.UserHandler.Register)
			protectedApi.POST("/send-notification", a.NotificationHandler.SendNotification)
			protectedApi.POST("/logout", a.UserHandler.Logout)
			protectedApi.POST("/gm-update", a.UserHandler.GmUpdate)
			protectedApi.POST("/normal-update", a.UserHandler.NormalUpdate)
			protectedApi.POST("/delete/:id", a.AuthMiddleware.RequireRole("gm"), a.UserHandler.GmDelete)
			// Other protected routes
		}
	}
}
