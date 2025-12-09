package leave

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveType string

const (
	fullDayLeave      LeaveType = "full day leave"
	halfDayLeave      LeaveType = "half day leave"
	unAuthorizedLeave LeaveType = "unauthorized leave"
)

type Status string

const (
	pending  Status = "pending"
	approved Status = "approved"
	rejected Status = "rejected"
)

type Leave struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	EmpID     primitive.ObjectID `bson:"emp_id" json:"emp_id"`
	LeaveType LeaveType          `bson:"leave_type" json:"leave_type"`
	StartDate primitive.DateTime `bson:"start_date" json:"start_date"`
	EndDate   primitive.DateTime `bson:"end_date" json:"end_date"`
	Reason    string             `bson:"reason" json:"reason"`
	Status    Status             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
