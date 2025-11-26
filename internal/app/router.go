package app

import "github.com/gin-gonic/gin"

// RegisterRoutes registers a minimal set of routes for the application.
func RegisterRoutes(r *gin.Engine, a *App) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Findix FMIIS running"})
	})

	api := r.Group("/api/v1")
	{
		// PUBLIC routes
		api.POST("/register", a.UserHandler.Register)
		api.POST("/login", a.UserHandler.Login)

		// PROTECTED routes
		protected := api.Group("")
		protected.Use(a.AuthMiddleware.Handle())
		{
			protected.POST("/send-notification", a.NotificationHandler.SendNotification)
			// Other protected routes
		}
	}
}
