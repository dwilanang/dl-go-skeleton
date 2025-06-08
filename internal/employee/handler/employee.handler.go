package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auth/util"
	"github.com/dwilanang/psp/internal/employee"
	"github.com/dwilanang/psp/internal/employee/dto"
	"github.com/dwilanang/psp/utils"
	utilrequest "github.com/dwilanang/psp/utils/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Deps employee.Dependencies
}

func NewHandler(deps employee.Dependencies) *Handler {
	return &Handler{
		Deps: deps,
	}
}

// CreateAttendance godoc
// @Security BearerAuth
// @Summary      Create attendance employee
// @Description  Create a new attendance employee
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        body  body      dto.AttendanceRequest  true  "attendance create payload"
// @Success      201   {object}  dto.EmployeeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/attendance [post]
func (h *Handler) CreateAttendance(c *gin.Context) {
	var ar dto.AttendanceRequest
	if err := c.ShouldBindJSON(&ar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	ar.By = userID
	ar.UserID = userID

	if err := h.Deps.Service.CreateAttendance(&ar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ar)
}

// FetchAttendance godoc
// @Security BearerAuth
// @Summary      List of attendacne
// @Description  Retrieves a list of attendacne with pagination.
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.AttendanceResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/attendance/all [get]
func (h *Handler) FetchAttendance(c *gin.Context) {
	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetAttendance(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get attendance"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateOvertime godoc
// @Security BearerAuth
// @Summary      Create overtime employee
// @Description  Create a new overtime employee
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        body  body      dto.OvertimeRequest  true  "overtime create payload"
// @Success      201   {object}  dto.EmployeeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/overtime [post]
func (h *Handler) CreateOvertime(c *gin.Context) {
	var or dto.OvertimeRequest
	if err := c.ShouldBindJSON(&or); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	or.By = userID
	or.UserID = userID

	if err := h.Deps.Service.CreateOvertime(&or); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, or)
}

// FetchOvertime godoc
// @Security BearerAuth
// @Summary      List of overtime
// @Description  Retrieves a list of overtime with pagination.
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.OvertimeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/overtime/all [get]
func (h *Handler) FetchOvertime(c *gin.Context) {
	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetOvertime(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get overtime"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateReimbursement godoc
// @Security BearerAuth
// @Summary      Create reimbursement employee
// @Description  Create a new reimbursement employee
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ReimbursementRequest  true  "reimbursement create payload"
// @Success      201   {object}  dto.EmployeeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/reimbursement [post]
func (h *Handler) CreateReimbursement(c *gin.Context) {
	var rr dto.ReimbursementRequest
	if err := c.ShouldBindJSON(&rr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	rr.By = userID
	rr.UserID = userID

	if err := h.Deps.Service.CreateReimbursement(&rr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rr)
}

// FetchReimbursement godoc
// @Security BearerAuth
// @Summary      List of reimbursement
// @Description  Retrieves a list of reimbursement with pagination.
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.ReimbursementResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/reimbursement/all [get]
func (h *Handler) FetchReimbursement(c *gin.Context) {
	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetReimbursement(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get reimbursement"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetEmployeePayslip godoc
// @Security BearerAuth
// @Summary      Get employee payslip
// @Description  View detailed payslip for an employee for a specific payroll
// @Tags         employee
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Payroll ID"
// @Success      201   {object}  dto.EmployeeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /employee/payslip/{id} [get]
func (h *Handler) GetEmployeePayslip(c *gin.Context) {
	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}

	id := c.Param("id")
	payrollID := utils.ConvertStringToInt(id)

	payslip, err := h.Deps.Service.GenerateEmployeePayslip(userID, payrollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payslip)
}
