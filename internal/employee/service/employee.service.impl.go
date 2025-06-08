package service

import (
	"errors"
	"fmt"
	"time"

	adminrepository "github.com/dwilanang/psp/internal/admin/repository"
	"github.com/dwilanang/psp/internal/employee/dto"
	"github.com/dwilanang/psp/internal/employee/model"
	"github.com/dwilanang/psp/internal/employee/repository"
	"github.com/dwilanang/psp/internal/employee/util"
	"github.com/dwilanang/psp/utils"
)

type service struct {
	repo      repository.Repository
	repoadmin adminrepository.Repository
}

func NewService(r repository.Repository, ra adminrepository.Repository) *service {
	return &service{repo: r, repoadmin: ra}
}

// CreateAttendance implements the Service interface.
func (s *service) CreateAttendance(request *dto.AttendanceRequest) error {
	// Parse date
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		fmt.Println("time.Parse() error: ", err)
		return errors.New("invalid date format")
	}

	periodID, err := s.getAttendancePeriodIDAndValidatePayroll(date)
	if err != nil {
		fmt.Println("s.getAttendancePeriodIDAndValidatePayroll() error: ", err)
		return err
	}

	// Submissions on non-working days.
	if !util.IsWeekday(date) {
		return errors.New("cannot submit attendance for non-working days (weekends or public holidays)")
	}

	if date.After(time.Now()) {
		return errors.New("cannot submit attendance for future date")
	}

	// Submissions on the same day should count as one.
	exists, err := s.repo.HasSubmittedAttendance(request.UserID, date)
	if err != nil {
		fmt.Println("s.repo.HasSubmittedAttendance() error: ", err)
		return err
	}
	if exists {
		fmt.Println("(s *service) CreateAttendance() error: ", err)
		return errors.New("attendance already submitted")
	}

	ap := &model.Attendance{
		UserID:         request.UserID,
		PeriodID:       periodID,
		DateAttendance: request.Date,
		CreatedBy:      request.By,
	}

	err = s.repo.CreateAttendance(ap)
	if err != nil {
		fmt.Println("s.repo.CreateAttendance() error: ", err)
	}

	return err
}

// GetAttendance implements the Service interface.
func (s *service) GetAttendance(userID int64, page, limit int) (*dto.AttendanceResponse, error) {
	var result dto.AttendanceResponse
	totalRecords, err := s.repo.CountAttendances(userID)
	if err != nil {
		fmt.Println("s.repo.CountAttendances() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.AttendanceData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.AttendanceData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	attendances, err := s.repo.FetchAttendances(userID, limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchAttendances() error: ", err)
		return &result, err
	}

	if attendances == nil {
		result.Data = []dto.AttendanceData{}
		return &result, nil
	}

	for _, v := range attendances {
		result.Data = append(result.Data, dto.AttendanceData{
			ID:           v.ID,
			FullName:     v.FullName,
			AttendanceAt: v.AttendanceAt,
		})
	}

	return &result, nil
}

// CreateOvertime implements the Service interface.
func (s *service) CreateOvertime(request *dto.OvertimeRequest) error {

	if request.Hours <= 0 {
		return errors.New("hours must be greater than 0")
	}

	// Overtime cannot be more than 3 hours per day.
	if request.Hours > 3 {
		return errors.New("hours must not exceed 3 per day")
	}

	// Parse date
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		fmt.Println("time.Parse() error: ", err)
		return errors.New("invalid date format")
	}

	periodID, err := s.getAttendancePeriodIDAndValidatePayroll(date)
	if err != nil {
		fmt.Println("s.getAttendancePeriodIDAndValidatePayroll() error: ", err)
		return err
	}

	// Overtime must be proposed after they are done working.
	now := time.Now()
	if date.Year() == now.Year() && date.YearDay() == now.YearDay() {
		limit := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, now.Location())
		if now.Before(limit) {
			return errors.New("overtime can only be submitted after working hours")
		}
	}

	// Submissions on the same day should count as one.
	exists, err := s.repo.HasSubmittedOvertime(request.UserID, date)
	if err != nil {
		fmt.Println("s.repo.HasSubmittedOvertime() error: ", err)
		return err
	}
	if exists {
		fmt.Println("(s *service) CreateOvertime() error: ", err)
		return errors.New("overtime already submitted")
	}

	ot := &model.Overtime{
		UserID:       request.UserID,
		PeriodID:     periodID,
		DateOvertime: request.Date,
		Hours:        request.Hours,
		CreatedBy:    request.By,
	}

	err = s.repo.CreateOvertimes(ot)
	if err != nil {
		fmt.Println("s.repo.CreateOvertimes() error: ", err)
	}

	return err
}

// GetOvertime implements the Service interface.
func (s *service) GetOvertime(userID int64, page, limit int) (*dto.OvertimeResponse, error) {
	var result dto.OvertimeResponse
	totalRecords, err := s.repo.CountOvertimes(userID)
	if err != nil {
		fmt.Println("s.repo.CountOvertimes() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.OvertimeData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.OvertimeData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	overtimes, err := s.repo.FetchOvertimes(userID, limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchOvertimes() error: ", err)
		return &result, err
	}

	if overtimes == nil {
		result.Data = []dto.OvertimeData{}
		return &result, nil
	}

	for _, v := range overtimes {
		result.Data = append(result.Data, dto.OvertimeData{
			ID:         v.ID,
			FullName:   v.FullName,
			OvertimeAt: v.OvertimeAt,
		})
	}

	return &result, nil
}

// CreateReimbursement implements the Service interface.
func (s *service) CreateReimbursement(request *dto.ReimbursementRequest) error {
	// Parse date
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		fmt.Println("time.Parse() error: ", err)
		return errors.New("invalid date format")
	}

	periodID, err := s.getAttendancePeriodIDAndValidatePayroll(date)
	if err != nil {
		fmt.Println("s.getAttendancePeriodIDAndValidatePayroll() error: ", err)
		return err
	}

	mr := &model.Reimbursement{
		UserID:      request.UserID,
		PeriodID:    periodID,
		Amount:      request.Amount,
		Description: request.Description,
		CreatedBy:   request.By,
	}

	err = s.repo.CreateReimbursement(mr)
	if err != nil {
		fmt.Println("s.repo.CreateAttendance() error: ", err)
	}

	return err
}

// GetReimbursement implements the Service interface.
func (s *service) GetReimbursement(userID int64, page, limit int) (*dto.ReimbursementResponse, error) {
	var result dto.ReimbursementResponse
	totalRecords, err := s.repo.CountReimbursement(userID)
	if err != nil {
		fmt.Println("s.repo.CountReimbursement() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.ReimbursementData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.ReimbursementData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	reimbersements, err := s.repo.FetchReimbursement(userID, limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchReimbursement() error: ", err)
		return &result, err
	}

	if reimbersements == nil {
		result.Data = []dto.ReimbursementData{}
		return &result, nil
	}

	for _, v := range reimbersements {
		result.Data = append(result.Data, dto.ReimbursementData{
			ID:          v.ID,
			FullName:    v.FullName,
			Amout:       v.Amount,
			Description: v.Description,
		})
	}

	return &result, nil
}

func (s *service) getAttendancePeriodIDAndValidatePayroll(date time.Time) (int64, error) {
	periodID, err := s.repoadmin.GetAttendancePeriodID(date)
	if err != nil {
		fmt.Println("s.repoadmin.GetAttendancePeriodID() error: ", err)
		return 0, err
	}
	if periodID == 0 {
		return 0, errors.New("no attendance period found for this date")
	}

	status, err := s.repoadmin.ValidatePayroll(periodID)
	if err != nil {
		fmt.Println("s.repoadmin.ValidatePayroll() error: ", err)
		return periodID, err
	}
	if status == "processed" {
		return periodID, errors.New("payroll already processed for this period")
	}

	return periodID, nil
}

func (s *service) GenerateEmployeePayslip(userID, payrollID int64) (*dto.EmployeePayslip, error) {
	periodID, err := s.repo.GetPayrollByID(payrollID)
	if err != nil {
		fmt.Println("s.repo.GetPayrollByID() error: ", err)
		return &dto.EmployeePayslip{}, err
	}

	employee, err := s.repo.GetEmployeeByID(userID)
	if err != nil {
		fmt.Println("s.repo.GetEmployeeByID() error: ", err)
		return &dto.EmployeePayslip{}, err
	}

	attendanceDays, err := s.repo.CalculateAttendances(userID, periodID)
	if err != nil {
		fmt.Println("s.repo.CalculateAttendances() error: ", err)
		return &dto.EmployeePayslip{}, err
	}

	totalOvertimeHours, err := s.repo.CalculateOvertimes(userID, periodID)
	if err != nil {
		fmt.Println("s.repo.CalculateOvertimes() error: ", err)
		return &dto.EmployeePayslip{}, err
	}
	overtimePay, err := s.repo.CalculateOvertimePay(userID, periodID, float64(employee.BaseSalary))
	if err != nil {
		fmt.Println("s.repo.GetEmployeeByID() error: ", err)
		return &dto.EmployeePayslip{}, err
	}

	totalReimburse, err := s.repo.CalculateReimburse(userID, periodID)
	if err != nil {
		fmt.Println("s.repo.GetEmployeeByID() error: ", err)
		return &dto.EmployeePayslip{}, err
	}

	attPay := float64(employee.BaseSalary) / 20.0 * float64(attendanceDays)
	thp := attPay + overtimePay + totalReimburse

	return &dto.EmployeePayslip{
		PayrollID:      payrollID,
		FullName:       employee.FullName,
		BaseSalary:     utils.RoundFloat(employee.BaseSalary, 2),
		AttendanceDays: attendanceDays,
		AttendancePay:  utils.RoundFloat(attPay, 2),
		OvertimeHours:  utils.RoundFloat(totalOvertimeHours, 0),
		OvertimePay:    utils.RoundFloat(overtimePay, 2),
		Reimbursements: utils.RoundFloat(totalReimburse, 2),
		TakeHomePay:    utils.RoundFloat(thp, 2),
	}, nil
}
