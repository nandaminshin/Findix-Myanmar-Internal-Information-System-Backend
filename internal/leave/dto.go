package leave

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaveRequest struct {
	EmpID     primitive.ObjectID `json:"emp_id" binding:"required"`
	LeaveType LeaveType          `json:"leave_type" binding:"required"`
	StartDate primitive.DateTime `json:"start_date" binding:"required"`
	EndDate   primitive.DateTime `json:"end_date" binding:"required"`
	Reason    string             `json:"reason" binding:"required"`
	Status    Status             `json:"status" binding:"required"`
}

type LeaveRespost struct {
	EmpID     primitive.ObjectID `json:"emp_id"`
	LeaveType LeaveType          `json:"leave_type"`
	StartDate primitive.DateTime `json:"start_date"`
	EndDate   primitive.DateTime `json:"end_date"`
	Reason    string             `json:"reason"`
	Status    Status             `json:"status"`
}
