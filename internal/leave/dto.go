package leave

type LeaveRequest struct {
	EmpID     string    `json:"emp_id" binding:"required"`
	LeaveType LeaveType `json:"leave_type" binding:"required"`
	StartDate string    `json:"start_date" binding:"required"`
	EndDate   string    `json:"end_date" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
	Status    Status    `json:"status" binding:"required"`
}

type LeaveResponse struct {
	ID        string    `json:"id"`
	EmpID     string    `json:"emp_id"`
	LeaveType LeaveType `json:"leave_type"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	Reason    string    `json:"reason"`
	Status    Status    `json:"status"`
}

type GmApprovalRequest struct {
	ID     string `json:"id" binding:"required"`
	Status Status `json:"status" binding:"required"`
}

type GmApprovalRespnse struct {
	ID        string    `json:"id"`
	EmpID     string    `json:"emp_id"`
	LeaveType LeaveType `json:"leave_type"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	Reason    string    `json:"reason"`
	Status    Status    `json:"status"`
	UpdatedAt string    `json:"updated_at"`
}
