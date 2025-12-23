package notification

import (
	"time"

	"fmiis/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotiType string

const (
	morningMeetingNnoti  NotiType = "Morning Meeting"
	devMeetingNoti       NotiType = "Developer Meeting"
	kosugiMeeting        NotiType = "Kosugi Meeting"
	emergencyMeetingNoti NotiType = "Emergency Meeting"
	internalMeetingNoti  NotiType = "Internal Meeting"
	generalNoti          NotiType = "General"
)

type Sender struct {
	Name  string    `bson:"name" json:"name"`
	Email string    `bson:"email" json:"email"`
	Image string    `bson:"image" json:"image"`
	Role  user.Role `bson:"role" json:"role"`
	Phone string    `bson:"phone" json:"phone"`
}

type Receiver struct {
	Name   string `bson:"name" json:"name"`
	Email  string `bson:"email" json:"email"`
	IsSeen bool   `bson:"is_seen" json:"is_seen"`
}

type Notification struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	NotiType  NotiType           `bson:"noti_type" json:"noti_type"`
	Sender    Sender             `bson:"sender" json:"sender"`
	Receivers []Receiver         `bson:"receivers" json:"receivers"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
