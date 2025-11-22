package app

import (
	"fmiis/internal/config"
	"fmiis/internal/database"
	"fmiis/internal/auth"
	"fmiis/internal/user"
	"fmiis/internal/middleware"
)

type App struct {
	Config      *config.Config
	DB          *database.MongoInstance
	AuthHandler *auth.Handler
	UserHandler *user.Handler
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
	// TODO: in future fill services/repo; for now create minimal handlers
	a.AuthHandler = NewAuthHandler(a)
	a.UserHandler = NewUserHandler(a)
}

// NewAuthHandler creates a minimal auth handler when other modules are not ready.
func NewAuthHandler(a *App) *auth.Handler {
	return &auth.Handler{}
}

// NewUserHandler creates a minimal user handler when other modules are not ready.
func NewUserHandler(a *App) *user.Handler {
	return &user.Handler{}
}

// StartServer is implemented in internal/app/server.go
