package notification

import (
	"time"
)

type NotificationRequest struct {
	NotiType  NotiType   `json:"noti_type" binding:"required"`
	Sender    Sender     `json:"sender" binding:"required"`
	Receivers []Receiver `json:"receivers" binding:"required"`
	Content   string     `json:"content" binding:"required"`
}

type NotificationResponse struct {
	ID        string     `json:"id"`
	NotiType  NotiType   `json:"noti_type"`
	Sender    Sender     `json:"sender"`
	Receivers []Receiver `json:"receivers"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
