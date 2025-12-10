package attendance

import "go.mongodb.org/mongo-driver/bson/primitive"

type AttendanceRequest struct {
	EmpID            string              `json:"emp_id" binding:"required"`
	Date             string              `json:"date" binding:"required"`
	AttendanceStatus AttendanceStatus    `json:"attendance_status" binding:"required"`
	LeaveID          *primitive.ObjectID `json:"leave_id,omitempty"`
}

type AttendanceResponse struct {
	ID               string           `json:"id"`
	EmpID            string           `json:"emp_id"`
	Date             string           `json:"date"`
	AttendanceStatus AttendanceStatus `json:"attendance_status"`
	LeaveID          string           `json:"leave_id,omitempty"`
}
