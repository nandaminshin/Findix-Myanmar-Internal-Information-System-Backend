package attendance

import "go.mongodb.org/mongo-driver/bson/primitive"

type AttendanceRequest struct {
	EmpID            primitive.ObjectID  `json:"emp_id" binding:"required"`
	Date             primitive.DateTime  `json:"date" binding:"required"`
	AttendanceStatus AttendanceStatus    `json:"attendance_status" binding:"required"`
	LeaveID          *primitive.ObjectID `json:"leave_id,omitempty"`
}

type AttendanceResponse struct {
	EmpID            primitive.ObjectID  `json:"emp_id"`
	Date             primitive.DateTime  `json:"date"`
	AttendanceStatus AttendanceStatus    `json:"attendance_status"`
	LeaveID          *primitive.ObjectID `json:"leave_id,omitempty"`
}
