package service

import (
	"fmt"

	"github.com/dwilanang/psp/internal/admin/dto"
	"github.com/dwilanang/psp/internal/admin/model"
	"github.com/dwilanang/psp/internal/admin/repository"
	"github.com/dwilanang/psp/utils"
)

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) *service {
	return &service{repo: r}
}

// CreateAttendancePeriods implements the Service interface.
func (s *service) CreateAttendancePeriods(request *dto.AttendancePeriodRequest) error {

	err := s.repo.ValidateAttendancePeriodDate(request.StartDate, request.EndDate)
	if err != nil {
		fmt.Println("s.repo.ValidateAttendancePeriodDate() error: ", err)
		return err
	}

	ap := &model.AttendancePeriod{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
		CreatedBy: request.By,
	}

	err = s.repo.CreateAttendancePeriods(ap)
	if err != nil {
		fmt.Println("s.repo.CreateAttendancePeriods() error: ", err)
	}

	return err
}

// GetAttendancePeriods implements the Service interface.
func (s *service) GetAttendancePeriods(page, limit int) (*dto.AttendancePeriodResponse, error) {
	var result dto.AttendancePeriodResponse
	totalRecords, err := s.repo.CountAttendancePeriods()
	if err != nil {
		fmt.Println("s.repo.CountAttendancePeriods() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.AttendancePeriodData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.AttendancePeriodData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	attendancePeriods, err := s.repo.FetchAttendancePeriods(limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchAttendancePeriods() error: ", err)
		return &result, err
	}

	if attendancePeriods == nil {
		result.Data = []dto.AttendancePeriodData{}
		return &result, nil
	}

	for _, a := range attendancePeriods {
		result.Data = append(result.Data, dto.AttendancePeriodData{
			ID:        a.ID,
			StartDate: a.StartDate,
			EndDate:   a.EndDate,
		})
	}

	return &result, nil
}

// CreatePayrolls implements the Service interface.
func (s *service) CreatePayrolls(request *dto.PayrollRequest) error {

	pr := &model.Payroll{
		PeriodID:  request.PeriodID,
		Status:    "pending",
		CreatedBy: request.By,
	}

	err := s.repo.CreatePayrolls(pr)
	if err != nil {
		fmt.Println("s.repo.CreatePayrolls() error: ", err)
	}

	return err
}

// GetPayrolls implements the Service interface.
func (s *service) GetPayrolls(page, limit int) (*dto.PayrollsResponse, error) {
	var result dto.PayrollsResponse
	totalRecords, err := s.repo.CountPayrolls()
	if err != nil {
		fmt.Println("s.repo.CountPayrolls() error: ", err)
		return &result, err
	}

	result.TotalRecord = totalRecords
	result.Page = page
	result.Limit = limit
	result.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if page > int(result.TotalPages) {
		result.Data = []dto.PayrollsResponseData{}
		return &result, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		result.Data = []dto.PayrollsResponseData{}
		return &result, nil
	}
	offset := (page - 1) * limit

	payrolls, err := s.repo.FetchPayrolls(limit, offset)
	if err != nil {
		fmt.Println("s.repo.FetchPayrolls() error: ", err)
		return &result, err
	}

	if payrolls == nil {
		result.Data = []dto.PayrollsResponseData{}
		return &result, nil
	}

	for _, v := range payrolls {
		result.Data = append(result.Data, dto.PayrollsResponseData{
			ID:          v.ID,
			Status:      v.Status,
			StartDate:   v.StartDate,
			EndDate:     v.EndDate,
			ProcessedAt: v.ProcessedAt,
		})
	}

	return &result, nil
}

// RunPayrolls implements the Service interface.
func (s *service) RunPayrolls(request *dto.PayrollRequest) error {

	pr := &model.Payroll{
		ID:        request.ID,
		Status:    "processed",
		UpdatedBy: request.By,
	}

	err := s.repo.UpdatePayrolls(pr)
	if err != nil {
		fmt.Println("s.repo.RunPayrolls() error: ", err)
	}

	return err
}

// SummaryPayrolls implements the Service interface.
func (s *service) SummaryPayrolls(request *dto.SummaryPayrollRequest) (*dto.PayrollSummaryResponse, error) {
	var summary dto.PayrollSummaryResponse
	periodID, err := s.repo.FindPayrollsByID(request.PayrollID)
	if err != nil {
		fmt.Println("s.repo.FindPayrollsByID() error: ", err)
		return &summary, err
	}

	periodInfo, err := s.repo.FindAttendancePeriodsByID(periodID)
	if err != nil {
		fmt.Println("s.repo.FindAttendancePeriodsByID() error: ", err)
		return &summary, err
	}

	if periodInfo == nil {
		return &summary, nil
	}

	totalRecords, err := s.repo.CountEmployeeByPeriodID(periodID)
	if err != nil {
		fmt.Println("s.repo.CountEmployeeByPeriodID() error: ", err)
		return &summary, err
	}

	summary.Period = dto.PeriodInfo{
		StartDate: periodInfo.StartDate,
		EndDate:   periodInfo.EndDate,
	}

	page := request.Page
	limit := request.Limit

	summary.TotalRecord = totalRecords
	summary.CurrentPage = page
	summary.Limit = limit
	summary.TotalPages = (totalRecords + int64(limit) - 1) / int64(limit)

	if request.Page > int(summary.TotalPages) {
		summary.Employees = []dto.EmployeePayslip{}
		return &summary, nil
	}

	if int64((page-1)*limit) >= totalRecords {
		summary.Employees = []dto.EmployeePayslip{}
		return &summary, nil
	}
	offset := (page - 1) * limit

	employees, err := s.repo.FindEmployeeByPeriodID(periodID, int64(request.Limit), int64(offset))
	if err != nil {
		fmt.Println("s.repo.FindEmployeeByPeriodID() error: ", err)
		return &summary, err
	}

	if employees == nil {
		return &summary, nil
	}

	summary.PayrollID = request.PayrollID

	total := float64(0)
	for _, r := range employees {
		attPay := float64(r.BaseSalary) / 20.0 * float64(r.AttendanceDays)
		thp := attPay + r.OvertimePay + r.Reimbursements

		summary.Employees = append(summary.Employees, dto.EmployeePayslip{
			ID:             r.ID,
			FullName:       r.FullName,
			BaseSalary:     utils.RoundFloat(r.BaseSalary, 2),
			AttendanceDays: r.AttendanceDays,
			AttendancePay:  utils.RoundFloat(attPay, 2),
			OvertimeHours:  utils.RoundFloat(r.OvertimeHours, 0),
			OvertimePay:    utils.RoundFloat(r.OvertimePay, 2),
			Reimbursements: utils.RoundFloat(r.Reimbursements, 2),
			TakeHomePay:    utils.RoundFloat(thp, 2),
		})
		total += thp
	}

	summary.TotalTakeHomePay = utils.RoundFloat(total, 2)

	return &summary, err
}
