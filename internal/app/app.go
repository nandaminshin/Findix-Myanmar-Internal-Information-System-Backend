package app

import (
	"fmiis/internal/auth"
	"fmiis/internal/common"
	"fmiis/internal/config"
	"fmiis/internal/database"
	"fmiis/internal/middleware"
	"fmiis/internal/notification"
	"fmiis/internal/user"
)

type App struct {
	Config              *config.Config
	DB                  *database.MongoInstance
	AuthHandler         *auth.Handler
	UserHandler         *user.UserHandler
	AuthMiddleware      middleware.AuthMiddleware
	NotificationHandler *notification.NotificationHandler
	Utilities           *common.Utilities
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := database.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		return nil, err
	}

	a := &App{
		Config: cfg,
		DB:     db,
	}
	a.initModules()
	return a, nil
}

func (a *App) initModules() {
	// Initialize User module
	userRepo := user.NewUserRepository(a.DB.DB)
	authService := auth.NewAuthService()
	userService := user.NewUserService(userRepo, authService)
	a.Utilities = common.NewUtility(userService)
	a.UserHandler = user.NewUserHandler(userService)
	a.AuthMiddleware = middleware.NewAuthMiddleware(authService)

	notificationRepo := notification.NewNotificationRepository(a.DB.DB)
	notificationService := notification.NewNotificationService(notificationRepo, *a.Utilities)
	a.NotificationHandler = notification.NewNotificationHandler(notificationService, userService)

	// Placeholder for AuthHandler if needed separately, or remove if merged
	a.AuthHandler = &auth.Handler{}

	//Utilities
}

// NewAuthHandler creates a minimal auth handler when other modules are not ready.
func NewAuthHandler(a *App) *auth.Handler {
	return &auth.Handler{}
}

// StartServer is implemented in internal/app/server.go
