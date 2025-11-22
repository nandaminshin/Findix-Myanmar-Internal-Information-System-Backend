package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func (a *App) StartServer() {
	// set Gin to release mode for production-like behavior
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// register application routes
	RegisterRoutes(r, a)

	// determine listen address
	port := "8000"
	if a != nil && a.Config != nil && a.Config.Port != "" {
		port = a.Config.Port
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("🚀 Findix server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}
}
