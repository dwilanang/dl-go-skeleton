package service

import (
	"errors"
	"testing"
	"time"

	repositoryadminmocks "github.com/dwilanang/psp/internal/admin/repository/mocks"
	"github.com/dwilanang/psp/internal/employee/dto"
	"github.com/dwilanang/psp/internal/employee/model"
	repositorymocks "github.com/dwilanang/psp/internal/employee/repository/mocks"
	"github.com/dwilanang/psp/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAttendance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	mockAdminRepo := repositoryadminmocks.NewMockRepository(ctrl)

	s := NewService(mockRepo, mockAdminRepo)

	userID := int64(1)
	by := int64(1)
	validDate := time.Now().AddDate(0, 0, -4).Format("2006-01-02") // yesterday
	weekendDate := "2025-06-07"                                    // Saturday

	futureDate := time.Now().AddDate(0, 0, 2).Format("2006-01-02") // 2 days later

	t.Run("success_create_attendance", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().HasSubmittedAttendance(userID, gomock.Any()).Return(false, nil)
		mockRepo.EXPECT().CreateAttendance(gomock.Any()).Return(nil)

		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   validDate,
			By:     by,
		})

		assert.NoError(t, err)
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   "2023-99-99",
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "invalid date format", err.Error())
	})

	t.Run("submit_on_weekend", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)

		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   weekendDate,
			By:     by,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot submit attendance for non-working days")
	})

	t.Run("submit_on_future_date", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)

		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   futureDate,
			By:     by,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot submit attendance for future date")
	})

	t.Run("attendance_already_submitted", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().HasSubmittedAttendance(userID, gomock.Any()).Return(true, nil)

		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   validDate,
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "attendance already submitted", err.Error())
	})

	t.Run("error_from_has_submitted_attendance", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().HasSubmittedAttendance(userID, gomock.Any()).Return(false, errors.New("db error"))

		err := s.CreateAttendance(&dto.AttendanceRequest{
			UserID: userID,
			Date:   validDate,
			By:     by,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestCreateOvertime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	mockAdminRepo := repositoryadminmocks.NewMockRepository(ctrl)

	s := NewService(mockRepo, mockAdminRepo)

	userID := int64(1)
	by := int64(1)
	validDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	t.Run("success_create_overtime", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().HasSubmittedOvertime(userID, gomock.Any()).Return(false, nil)
		mockRepo.EXPECT().CreateOvertimes(gomock.Any()).Return(nil)

		err := s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   validDate,
			Hours:  2,
			By:     by,
		})
		assert.NoError(t, err)
	})

	t.Run("hours_less_or_equal_zero", func(t *testing.T) {
		err := s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   validDate,
			Hours:  0,
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "hours must be greater than 0", err.Error())
	})

	t.Run("hours_exceeds_limit", func(t *testing.T) {
		err := s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   validDate,
			Hours:  4,
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "hours must not exceed 3 per day", err.Error())
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		err := s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   "2023-99-99",
			Hours:  2,
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "invalid date format", err.Error())
	})

	t.Run("overtime_before_5pm_today", func(t *testing.T) {
		today := "2023-07-08 18:00:00"

		_ = s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   today,
			Hours:  1,
			By:     by,
		})
	})

	t.Run("overtime_already_submitted", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().HasSubmittedOvertime(userID, gomock.Any()).Return(true, nil)

		err := s.CreateOvertime(&dto.OvertimeRequest{
			UserID: userID,
			Date:   validDate,
			Hours:  2,
			By:     by,
		})
		assert.Error(t, err)
		assert.Equal(t, "overtime already submitted", err.Error())
	})
}

func TestCreateReimbursement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	mockAdminRepo := repositoryadminmocks.NewMockRepository(ctrl)

	s := NewService(mockRepo, mockAdminRepo)

	userID := int64(1)
	by := int64(1)
	validDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	t.Run("success_create_reimbursement", func(t *testing.T) {
		mockAdminRepo.EXPECT().GetAttendancePeriodID(gomock.Any()).Return(int64(10), nil)
		mockAdminRepo.EXPECT().ValidatePayroll(int64(10)).Return("", nil)
		mockRepo.EXPECT().CreateReimbursement(gomock.Any()).Return(nil)

		err := s.CreateReimbursement(&dto.ReimbursementRequest{
			UserID:      userID,
			Date:        validDate,
			Amount:      100,
			Description: "Taxi fare",
			By:          by,
		})
		assert.NoError(t, err)
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		err := s.CreateReimbursement(&dto.ReimbursementRequest{
			UserID:      userID,
			Date:        "2023-99-99",
			Amount:      100,
			Description: "Desc",
			By:          by,
		})
		assert.Error(t, err)
		assert.Equal(t, "invalid date format", err.Error())
	})
}

func TestGenerateEmployeePayslip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	mockAdminRepo := repositoryadminmocks.NewMockRepository(ctrl)

	s := NewService(mockRepo, mockAdminRepo)

	userID := int64(123)
	payrollID := int64(456)

	employee := model.Employee{
		FullName:   "Pororo",
		BaseSalary: 1000000,
	}

	t.Run("success", func(t *testing.T) {
		periodID := int64(789)
		attendanceDays := int64(18)
		totalOvertimeHours := 5.0
		overtimePay := 50000.0
		totalReimburse := 75000.0

		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(employee, nil)
		mockRepo.EXPECT().CalculateAttendances(userID, periodID).Return(attendanceDays, nil)
		mockRepo.EXPECT().CalculateOvertimes(userID, periodID).Return(totalOvertimeHours, nil)
		mockRepo.EXPECT().CalculateOvertimePay(userID, periodID, float64(employee.BaseSalary)).Return(overtimePay, nil)
		mockRepo.EXPECT().CalculateReimburse(userID, periodID).Return(totalReimburse, nil)

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.NoError(t, err)
		assert.Equal(t, payrollID, payslip.PayrollID)
		assert.Equal(t, employee.FullName, payslip.FullName)
		assert.Equal(t, utils.RoundFloat(employee.BaseSalary, 2), payslip.BaseSalary)
		assert.Equal(t, attendanceDays, payslip.AttendanceDays)
		assert.Equal(t, utils.RoundFloat(float64(employee.BaseSalary)/20.0*float64(attendanceDays), 2), payslip.AttendancePay)
		assert.Equal(t, utils.RoundFloat(totalOvertimeHours, 0), payslip.OvertimeHours)
		assert.Equal(t, utils.RoundFloat(overtimePay, 2), payslip.OvertimePay)
		assert.Equal(t, utils.RoundFloat(totalReimburse, 2), payslip.Reimbursements)
		assert.Equal(t, utils.RoundFloat((float64(employee.BaseSalary)/20.0*float64(attendanceDays))+overtimePay+totalReimburse, 2), payslip.TakeHomePay)
	})

	t.Run("error_get_payroll", func(t *testing.T) {
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(int64(0), errors.New("payroll not found"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "payroll not found", err.Error())
		assert.Empty(t, payslip.FullName)
	})

	t.Run("error_get_employee", func(t *testing.T) {
		periodID := int64(789)
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(model.Employee{}, errors.New("employee not found"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "employee not found", err.Error())
		assert.Empty(t, payslip.FullName)
	})

	t.Run("error_calculate_attendance", func(t *testing.T) {
		periodID := int64(789)
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(employee, nil)
		mockRepo.EXPECT().CalculateAttendances(userID, periodID).Return(int64(0), errors.New("failed to calculate attendance"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "failed to calculate attendance", err.Error())
		assert.Empty(t, payslip.FullName)
	})

	t.Run("error_calculate_overtime", func(t *testing.T) {
		periodID := int64(789)
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(employee, nil)
		mockRepo.EXPECT().CalculateAttendances(userID, periodID).Return(int64(18), nil)
		mockRepo.EXPECT().CalculateOvertimes(userID, periodID).Return(float64(0), errors.New("failed to calculate overtime"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "failed to calculate overtime", err.Error())
		assert.Empty(t, payslip.FullName)
	})

	t.Run("error_calculate_overtime_pay", func(t *testing.T) {
		periodID := int64(789)
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(employee, nil)
		mockRepo.EXPECT().CalculateAttendances(userID, periodID).Return(int64(18), nil)
		mockRepo.EXPECT().CalculateOvertimes(userID, periodID).Return(float64(5.0), nil)
		mockRepo.EXPECT().CalculateOvertimePay(userID, periodID, float64(employee.BaseSalary)).Return(float64(0), errors.New("failed to calculate overtime pay"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "failed to calculate overtime pay", err.Error())
		assert.Empty(t, payslip.FullName)
	})

	t.Run("error_calculate_reimburse", func(t *testing.T) {
		periodID := int64(789)
		mockRepo.EXPECT().GetPayrollByID(payrollID).Return(periodID, nil)
		mockRepo.EXPECT().GetEmployeeByID(userID).Return(employee, nil)
		mockRepo.EXPECT().CalculateAttendances(userID, periodID).Return(int64(18), nil)
		mockRepo.EXPECT().CalculateOvertimes(userID, periodID).Return(float64(5.0), nil)
		mockRepo.EXPECT().CalculateOvertimePay(userID, periodID, float64(employee.BaseSalary)).Return(float64(50000.0), nil)
		mockRepo.EXPECT().CalculateReimburse(userID, periodID).Return(float64(0), errors.New("failed to calculate reimburse"))

		payslip, err := s.GenerateEmployeePayslip(userID, payrollID)
		assert.Error(t, err)
		assert.Equal(t, "failed to calculate reimburse", err.Error())
		assert.Empty(t, payslip.FullName)
	})
}
