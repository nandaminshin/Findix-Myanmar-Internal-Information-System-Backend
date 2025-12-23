package leave

import (
	"context"
	"errors"
	"fmiis/internal/attendance"
	"fmiis/internal/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveService interface {
	RecordLeaveRequest(ctx context.Context, req *LeaveRequest) (*LeaveResponse, error)
	GmApproval(ctx context.Context, req *GmApprovalRequest) (*GmApprovalRespnse, error)
}

type leaveService struct {
	repo           LeaveRepository
	utilities      common.Utilities
	attendanceRepo attendance.AttendanceRepository
}

func NewLeaveService(repo LeaveRepository, u common.Utilities, attRepo attendance.AttendanceRepository) LeaveService {
	return &leaveService{
		repo:           repo,
		utilities:      u,
		attendanceRepo: attRepo,
	}
}

func IsValidLeaveType(leaveType LeaveType) bool {
	switch leaveType {
	case morningMeetingLeave, developerMeetingLeave, globleTeamMeetingLeave, fullDayLeave, halfDayLeave, unAuthorizedLeave:
		return true
	default:
		return false
	}
}

func IsValidLeaveStatus(status Status) bool {
	switch status {
	case pending, approved, rejected:
		return true
	default:
		return false
	}
}

func (s *leaveService) RecordLeaveRequest(ctx context.Context, req *LeaveRequest) (*LeaveResponse, error) {
	if !IsValidLeaveType(req.LeaveType) {
		return nil, errors.New("invalid leave type")
	}
	if !IsValidLeaveStatus(req.Status) {
		return nil, errors.New("invalid leave status")
	}

	EmpObjID, err := s.utilities.ParseObjectIDForDB(req.EmpID)
	if err != nil {
		return nil, err
	}

	startDate, err := s.utilities.ParseDateTimeForDB(req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := s.utilities.ParseDateTimeForDB(req.EndDate)
	if err != nil {
		return nil, err
	}

	existingAttendance, err := s.repo.FindBymEmpIDAndEndDate(ctx, *EmpObjID, *startDate, *endDate)
	if err != nil {
		return nil, err
	}
	if existingAttendance != nil {
		return nil, errors.New("leave day already exists")
	}

	leave := &Leave{
		EmpID:     *EmpObjID,
		LeaveType: req.LeaveType,
		StartDate: *startDate,
		EndDate:   *endDate,
		Reason:    req.Reason,
		Status:    "pending",
	}

	dberr := s.repo.Create(ctx, leave)
	if dberr != nil {
		return nil, err
	}

	response := &LeaveResponse{
		ID:        leave.ID.Hex(),
		EmpID:     req.EmpID,
		LeaveType: req.LeaveType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Reason:    req.Reason,
		Status:    "pending",
	}

	return response, nil
}

func (s *leaveService) GmApproval(ctx context.Context, req *GmApprovalRequest) (*GmApprovalRespnse, error) {
	if !IsValidLeaveStatus(req.Status) {
		return nil, errors.New("invalid leave status")
	}

	leaveObjID, err := s.utilities.ParseObjectIDForDB(req.ID)
	if err != nil {
		return nil, err
	}
	leave, err := s.repo.FindByID(ctx, *leaveObjID)
	if err != nil {
		return nil, err
	}

	leave.Status = req.Status
	leave.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, leave)
	if err != nil {
		return nil, err
	}

	// If approved → create attendance records
	if req.Status == "approved" {

		start := leave.StartDate.Time()
		end := leave.EndDate.Time()

		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {

			att := &attendance.Attendance{
				ID:               primitive.NewObjectID(),
				Date:             primitive.NewDateTimeFromTime(d),
				AttendanceStatus: attendance.AttendanceStatus("absant"),
				LeaveID:          &leave.ID,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			if err := s.attendanceRepo.Create(ctx, att); err != nil {
				return nil, err
			}
		}
	}

	res := &GmApprovalRespnse{
		ID:        leave.ID.Hex(),
		Status:    leave.Status,
		UpdatedAt: leave.UpdatedAt.String(),
	}

	return res, nil
}
