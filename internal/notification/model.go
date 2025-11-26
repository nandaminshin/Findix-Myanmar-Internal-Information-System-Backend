package notification

import (
	"time"

	"fmiis/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotiType string

const (
	morning_meeting_noti   NotiType = "morning_meeting_noti"
	dev_meeting_noti       NotiType = "dev_meeting_noti"
	emergency_meeting_noti NotiType = "emergency_meeting_noti"
	general_noti           NotiType = "general_noti"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	NotiType  NotiType           `bson:"noti_type" json:"noti_type"`
	Sender    user.User          `bson:"sender" json:"sender"`
	Receiver  user.User          `bson:"receiver" json:"receiver"`
	Content   string             `bson:"content" json:"content"`
	IsSeen    bool               `bson:"is_seen" json:"is_seen"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
