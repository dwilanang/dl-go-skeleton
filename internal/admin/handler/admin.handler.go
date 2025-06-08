package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/admin"
	"github.com/dwilanang/psp/internal/admin/dto"
	"github.com/dwilanang/psp/internal/auth/util"
	"github.com/dwilanang/psp/utils"
	utilrequest "github.com/dwilanang/psp/utils/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Deps admin.Dependencies
}

func NewHandler(deps admin.Dependencies) *Handler {
	return &Handler{
		Deps: deps,
	}
}

// CreateAttendancePeriod godoc
// @Security BearerAuth
// @Summary      Create attendance period
// @Description  Create a new attendance period
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body  body      dto.AttendancePeriodRequest  true  "attendance period create payload"
// @Success      201   {object}  dto.AdminResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/attendance-periods [post]
func (h *Handler) CreateAttendancePeriod(c *gin.Context) {
	var ar dto.AttendancePeriodRequest
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

	if err := h.Deps.Service.CreateAttendancePeriods(&ar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create attendance period"})
		return
	}

	c.JSON(http.StatusCreated, ar)
}

// FetchAttendancePeriod godoc
// @Security BearerAuth
// @Summary      List of attendance period
// @Description  Retrieves a list of attendance period with pagination.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.AttendancePeriodResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/attendance-periods/all [get]
func (h *Handler) FetchAttendancePeriod(c *gin.Context) {

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetAttendancePeriods(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get attendance period"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreatePayroll godoc
// @Security BearerAuth
// @Summary      Create payroll
// @Description  Create a new payroll
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body  body      dto.PayrollRequest  true  "period create payload"
// @Success      201   {object}  dto.AdminResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/payroll [post]
func (h *Handler) CreatePayroll(c *gin.Context) {
	var pr dto.PayrollRequest
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	pr.By = userID

	if err := h.Deps.Service.CreatePayrolls(&pr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not run payroll"})
		return
	}

	c.JSON(http.StatusCreated, pr)
}

// FetchPayroll godoc
// @Security BearerAuth
// @Summary      List of payroll
// @Description  Retrieves a list of payroll with pagination.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.PayrollsResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/payroll/all [get]
func (h *Handler) FetchPayroll(c *gin.Context) {

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetPayrolls(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get payroll"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RunPayroll godoc
// @Security BearerAuth
// @Summary      Run payroll
// @Description  Run a new payroll
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Payroll ID"
// @Success      201   {object}  dto.AdminResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/payroll/{id}/process [put]
func (h *Handler) RunPayroll(c *gin.Context) {
	id := c.Param("id")

	var pp dto.PayrollRequest

	userID, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	pp.By = userID

	idInt := utils.ConvertStringToInt(id)
	pp.ID = idInt

	if err := h.Deps.Service.RunPayrolls(&pp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not run payroll"})
		return
	}

	c.JSON(http.StatusCreated, pp)
}

// SummaryPayroll godoc
// @Security BearerAuth
// @Summary      Summary payroll
// @Description  Summary a payroll
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id    path     int  true  "Payroll ID"
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.AdminResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /admin/payroll/{id}/summary [get]
func (h *Handler) SummaryPayroll(c *gin.Context) {
	id := c.Param("id")

	var spr dto.SummaryPayrollRequest

	idInt := utils.ConvertStringToInt(id)
	spr.PayrollID = idInt

	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	spr.Page = page
	spr.Limit = limit

	resp, err := h.Deps.Service.SummaryPayrolls(&spr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not getsummary"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
