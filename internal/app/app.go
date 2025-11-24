package app

import (
	"fmiis/internal/auth"
	"fmiis/internal/config"
	"fmiis/internal/database"
	"fmiis/internal/middleware"
	"fmiis/internal/user"
)

type App struct {
	Config      *config.Config
	DB          *database.MongoInstance
	AuthHandler *auth.Handler
	UserHandler *user.UserHandler
	Middlewares *middleware.Manager
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
	a.Middlewares = middleware.NewManager(cfg)
	a.initModules()
	return a, nil
}

func (a *App) initModules() {
	// Initialize User module
	userRepo := user.NewUserRepository(a.DB.DB)
	authService := auth.NewAuthService()
	userService := user.NewUserService(userRepo, authService)
	a.UserHandler = user.NewUserHandler(userService)

	// Placeholder for AuthHandler if needed separately, or remove if merged
	a.AuthHandler = &auth.Handler{}
}

// NewAuthHandler creates a minimal auth handler when other modules are not ready.
func NewAuthHandler(a *App) *auth.Handler {
	return &auth.Handler{}
}

// StartServer is implemented in internal/app/server.go
