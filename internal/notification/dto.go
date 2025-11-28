package notification

import (
	"fmiis/internal/user"
	"time"
)

type NotificationRequest struct {
	NotiType  NotiType   `json:"noti_type" binding:"required"`
	Sender    user.User  `json:"sender" binding:"required"`
	Receivers []Receiver `json:"receivers" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	IsSeen    bool       `json:"is_seen"`
}

type NotificationResponse struct {
	NotiType  NotiType   `json:"noti_type" binding:"required"`
	Sender    user.User  `json:"sender" binding:"required"`
	Receivers []Receiver `json:"receivers" binding:"required"`
	Content   string     `json:"content" binding:"required"`
	IsSeen    bool       `json:"is_seen"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
