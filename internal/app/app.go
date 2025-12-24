package app

import (
	"context"
	"fmiis/internal/attendance"
	"fmiis/internal/auth"
	"fmiis/internal/common"
	"fmiis/internal/config"
	"fmiis/internal/database"
	"fmiis/internal/leave"
	"fmiis/internal/middleware"
	"fmiis/internal/normal_email"
	"fmiis/internal/notification"
	"fmiis/internal/storage"
	"fmiis/internal/user"
	"log"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type App struct {
	Config              *config.Config
	DB                  *database.MongoInstance
	SocketServer        *socketio.Server
	StorageService      storage.StorageService
	AuthHandler         *auth.Handler
	UserHandler         *user.UserHandler
	AuthMiddleware      middleware.AuthMiddleware
	NotificationHandler *notification.NotificationHandler
	Utilities           *common.Utilities
	NormalEmailService  normal_email.EmailService
	AttendanceHandler   *attendance.AttendanceHandler
	LeaveHandler        *leave.LeaveHandler
}

func NewApp(cfg *config.Config, server socketio.Server) (*App, error) {
	db, err := database.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		return nil, err
	}

	storageService := storage.NewSupabaseStorage(cfg.SupabaseURL, cfg.SupabaseServiceKey, cfg.SupabaseBucket)

	a := &App{
		Config:         cfg,
		DB:             db,
		SocketServer:   &server,
		StorageService: storageService,
	}
	a.initModules()
	return a, nil
}

func (a *App) initModules() {
	// Initialize User module
	userRepo := user.NewUserRepository(a.DB.DB)
	authService := auth.NewAuthService()
	userService := user.NewUserService(userRepo, authService, a.SocketServer, a.StorageService)
	a.Utilities = common.NewUtility(userService)
	a.UserHandler = user.NewUserHandler(userService)
	a.AuthMiddleware = middleware.NewAuthMiddleware(authService)

	a.NormalEmailService = normal_email.NewBrevoService(a.Config.BrevoAPIKey, a.Config.EmailFrom, a.Config.EmailFromName)
	notificationRepo := notification.NewNotificationRepository(a.DB.DB)
	notificationService := notification.NewNotificationService(notificationRepo, *a.Utilities, a.NormalEmailService, userService)
	a.NotificationHandler = notification.NewNotificationHandler(notificationService, userService)
	setupTTLIndex(notificationRepo)

	attendanceRepo := attendance.NewAttendanceRepository(a.DB.DB)
	attendanceService := attendance.NewAttendanceService(attendanceRepo, *a.Utilities)
	a.AttendanceHandler = attendance.NewAttendanceHandler(attendanceService)

	leaveRepo := leave.NewLeaveRepository(a.DB.DB)
	leaveService := leave.NewLeaveService(leaveRepo, *a.Utilities, attendanceRepo)
	a.LeaveHandler = leave.NewLeaveHandler(leaveService)

	// Placeholder for AuthHandler if needed separately, or remove if merged
	a.AuthHandler = &auth.Handler{}

	//Utilities
}

func setupTTLIndex(repo notification.NotificationRepository) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := repo.SetupTTLIndex(ctx)
	if err != nil {
		// Check if it's just "index already exists" error
		if isIndexExistsError(err) {
			log.Println("✅ MongoDB TTL index already exists")
		} else {
			log.Printf("⚠️ Failed to setup TTL index: %v", err)
			log.Printf("ℹ️ Notifications will not be auto-deleted. Manual cleanup may be needed.")
		}
	} else {
		log.Println("✅ MongoDB TTL index created successfully")
		log.Println("📅 Notifications will be automatically deleted after 7 days")
	}
}

// Helper function to check if index already exists
func isIndexExistsError(err error) bool {
	// MongoDB returns specific error codes for duplicate indexes
	// You might need to adjust this based on your MongoDB driver version
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "index already exists") || strings.Contains(msg, "duplicate key")
}

// NewAuthHandler creates a minimal auth handler when other modules are not ready.
func NewAuthHandler(a *App) *auth.Handler {
	return &auth.Handler{}
}
