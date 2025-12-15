package attendance

import (
	"context"
	"errors"

	"fmiis/internal/common"
)

type AttendanceService interface {
	RecordAttendance(ctx context.Context, req *AttendanceRequest) (*AttendanceResponse, error)
}

type attendanceService struct {
	repo      AttendanceRepository
	utilities common.Utilities
}

func NewAttendanceService(repo AttendanceRepository, u common.Utilities) AttendanceService {
	return &attendanceService{
		repo:      repo,
		utilities: u,
	}
}

func (s *attendanceService) RecordAttendance(ctx context.Context, req *AttendanceRequest) (*AttendanceResponse, error) {
	EmpObjID, err := s.utilities.ParseObjectIDForDB(req.EmpID)
	if err != nil {
		return nil, err
	}

	date, err := s.utilities.ParseDateTimeForDB(req.Date)
	if err != nil {
		return nil, err
	}

	existingAttendance, err := s.repo.FindBymEmpIDAndDate(ctx, *EmpObjID, *date)
	if err != nil {
		return nil, err
	}
	if existingAttendance != nil {
		return nil, errors.New("attendance already exists")
	}

	attendance := &Attendance{
		EmpID:            *EmpObjID,
		Date:             *date,
		AttendanceStatus: req.AttendanceStatus,
	}
	if req.LeaveID != nil {
		attendance.LeaveID = req.LeaveID
	}

	if err := s.repo.Create(ctx, attendance); err != nil {
		return nil, err
	}

	response := &AttendanceResponse{
		ID:               attendance.ID.Hex(),
		EmpID:            attendance.EmpID.Hex(),
		Date:             attendance.Date.Time().Format("2006-01-02"),
		AttendanceStatus: attendance.AttendanceStatus,
	}

	if attendance.LeaveID != nil {
		response.LeaveID = attendance.LeaveID.Hex()
	}

	return response, nil
}
