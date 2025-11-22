package middleware

import "fmiis/internal/config"

// Manager is a placeholder for middleware dependencies.
type Manager struct {
    cfg *config.Config
}

// NewManager creates a new middleware manager (minimal implementation).
func NewManager(cfg *config.Config) *Manager {
    return &Manager{cfg: cfg}
}
