package app

import "github.com/gin-gonic/gin"

// RegisterRoutes registers a minimal set of routes for the application.
func RegisterRoutes(r *gin.Engine, a *App) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Findix FMIIS running"})
	})
}
