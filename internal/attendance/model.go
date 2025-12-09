package attendance

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttendanceStatus string

const (
	present AttendanceStatus = "present"
	absant  AttendanceStatus = "absant"
)

type Attendance struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	EmpID            primitive.ObjectID `bson:"emp_id" json:"emp_id"`
	Date             primitive.DateTime `bson:"date" json:"date"`
	AttendanceStatus AttendanceStatus   `bson:"attendance_status" json:"attendance_status"`
	LeaveID          primitive.ObjectID `bson:"leave_id" json:"leave_id"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
