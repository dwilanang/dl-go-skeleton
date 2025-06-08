package service

import (
	"errors"
	"testing"
	"time"

	"github.com/dwilanang/psp/internal/admin/dto"
	"github.com/dwilanang/psp/internal/admin/model"
	repositorymocks "github.com/dwilanang/psp/internal/admin/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAttendancePeriods_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.AttendancePeriodRequest{
		StartDate: "2024-06-01",
		EndDate:   "2024-06-30",
		By:        1,
	}

	mockRepo.EXPECT().
		ValidateAttendancePeriodDate(req.StartDate, req.EndDate).
		Return(nil)

	mockRepo.EXPECT().
		CreateAttendancePeriods(gomock.Any()).
		Return(nil)

	err := svc.CreateAttendancePeriods(req)
	assert.NoError(t, err)
}

func TestCreateAttendancePeriods_ValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.AttendancePeriodRequest{
		StartDate: "2024-06-01",
		EndDate:   "2024-06-30",
		By:        1,
	}

	mockRepo.EXPECT().
		ValidateAttendancePeriodDate(req.StartDate, req.EndDate).
		Return(errors.New("invalid date"))

	err := svc.CreateAttendancePeriods(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid date")
}

func TestCreateAttendancePeriods_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.AttendancePeriodRequest{
		StartDate: "2024-06-01",
		EndDate:   "2024-06-30",
		By:        1,
	}

	mockRepo.EXPECT().
		ValidateAttendancePeriodDate(req.StartDate, req.EndDate).
		Return(nil)

	mockRepo.EXPECT().
		CreateAttendancePeriods(gomock.Any()).
		Return(errors.New("insert error"))

	err := svc.CreateAttendancePeriods(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "insert error")
}

func TestCreatePayrolls_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.PayrollRequest{
		PeriodID: 123,
		By:       1,
	}

	mockRepo.EXPECT().
		CreatePayrolls(gomock.Any()).
		Return(nil)

	err := svc.CreatePayrolls(req)
	assert.NoError(t, err)
}

func TestCreatePayrolls_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.PayrollRequest{
		PeriodID: 123,
		By:       1,
	}

	mockRepo.EXPECT().
		CreatePayrolls(gomock.Any()).
		Return(errors.New("insert error"))

	err := svc.CreatePayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "insert error")
}

func TestRunPayrolls_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.PayrollRequest{
		ID: 1,
		By: 2,
	}

	mockRepo.EXPECT().
		UpdatePayrolls(gomock.Any()).
		Return(nil)

	err := svc.RunPayrolls(req)
	assert.NoError(t, err)
}

func TestRunPayrolls_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.PayrollRequest{
		ID: 1,
		By: 2,
	}

	mockRepo.EXPECT().
		UpdatePayrolls(gomock.Any()).
		Return(errors.New("update error"))

	err := svc.RunPayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "update error")
}

func TestService_SummaryPayrolls_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)

	svc := &service{repo: mockRepo}

	req := &dto.SummaryPayrollRequest{
		PayrollID: 5,
		Page:      1,
		Limit:     2,
	}

	periodID := int64(10)
	mockRepo.EXPECT().
		FindPayrollsByID(req.PayrollID).
		Return(periodID, nil)

	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	periodInfo := &model.PeriodInfo{
		ID:        periodID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	mockRepo.EXPECT().
		FindAttendancePeriodsByID(periodID).
		Return(periodInfo, nil)

	mockRepo.EXPECT().
		CountEmployeeByPeriodID(periodID).
		Return(int64(2), nil)

	employeeList := []*model.EmployeePayslip{
		{
			ID:             1,
			FullName:       "Alice",
			BaseSalary:     4000000,
			AttendanceDays: 20,
			OvertimeHours:  5,
			OvertimePay:    500000,
			Reimbursements: 200000,
		},
		{
			ID:             2,
			FullName:       "Bob",
			BaseSalary:     5000000,
			AttendanceDays: 18,
			OvertimeHours:  3,
			OvertimePay:    300000,
			Reimbursements: 100000,
		},
	}
	mockRepo.EXPECT().
		FindEmployeeByPeriodID(periodID, int64(req.Limit), int64(0)).
		Return(employeeList, nil)

	resp, err := svc.SummaryPayrolls(req)

	assert.NoError(t, err)
	assert.Equal(t, req.PayrollID, resp.PayrollID)
	assert.Equal(t, int64(2), resp.TotalRecord)
	assert.Equal(t, 1, resp.CurrentPage)
	assert.Equal(t, 2, resp.Limit)
	assert.Equal(t, int64(1), resp.TotalPages)
	assert.Equal(t, startDate, resp.Period.StartDate)
	assert.Equal(t, endDate, resp.Period.EndDate)
	assert.Len(t, resp.Employees, 2)
	assert.Greater(t, resp.TotalTakeHomePay, 0.0)
}

func TestSummaryPayrolls_FindPayrollsByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.SummaryPayrollRequest{
		PayrollID: 1,
		Page:      1,
		Limit:     2,
	}

	mockRepo.EXPECT().FindPayrollsByID(req.PayrollID).Return(int64(0), errors.New("not found"))

	summary, err := svc.SummaryPayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Equal(t, int64(0), summary.PayrollID)
}

func TestSummaryPayrolls_FindAttendancePeriodsByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.SummaryPayrollRequest{
		PayrollID: 1,
		Page:      1,
		Limit:     2,
	}

	mockRepo.EXPECT().FindPayrollsByID(req.PayrollID).Return(int64(10), nil)
	mockRepo.EXPECT().FindAttendancePeriodsByID(int64(10)).Return(nil, errors.New("period not found"))

	summary, err := svc.SummaryPayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "period not found")
	assert.Equal(t, int64(0), summary.PayrollID)
}

func TestSummaryPayrolls_CountEmployeeByPeriodID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.SummaryPayrollRequest{
		PayrollID: 1,
		Page:      1,
		Limit:     2,
	}

	periodID := int64(10)

	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	periodInfo := &model.PeriodInfo{
		ID:        periodID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	mockRepo.EXPECT().FindPayrollsByID(req.PayrollID).Return(int64(10), nil)
	mockRepo.EXPECT().
		FindAttendancePeriodsByID(periodID).
		Return(periodInfo, nil)
	mockRepo.EXPECT().CountEmployeeByPeriodID(int64(10)).Return(int64(0), errors.New("count error"))

	summary, err := svc.SummaryPayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "count error")
	assert.Equal(t, int64(0), summary.PayrollID)
}

func TestSummaryPayrolls_FindEmployeeByPeriodID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositorymocks.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.SummaryPayrollRequest{
		PayrollID: 1,
		Page:      1,
		Limit:     2,
	}

	periodID := int64(10)

	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	periodInfo := &model.PeriodInfo{
		ID:        periodID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	mockRepo.EXPECT().FindPayrollsByID(req.PayrollID).Return(int64(10), nil)
	mockRepo.EXPECT().
		FindAttendancePeriodsByID(periodID).
		Return(periodInfo, nil)
	mockRepo.EXPECT().CountEmployeeByPeriodID(int64(10)).Return(int64(2), nil)
	mockRepo.EXPECT().FindEmployeeByPeriodID(int64(10), int64(req.Limit), int64(0)).Return(nil, errors.New("employee error"))

	summary, err := svc.SummaryPayrolls(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "employee error")
	assert.Equal(t, int64(0), summary.PayrollID)
}
