package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auditlogs"
	"github.com/dwilanang/psp/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Deps auditlogs.Dependencies
}

func NewHandler(deps auditlogs.Dependencies) *Handler {
	return &Handler{
		Deps: deps,
	}
}

// GetAuditLogs godoc
// @Security BearerAuth
// @Summary      List of audit logs
// @Description  Retrieves a list of audit logs with pagination.
// @Tags         auditlogs
// @Accept       json
// @Produce      json
// @Param        page  query    int  false "Page number (default: 1)"
// @Param        limit query    int  false "Items per page (default: 20)"
// @Success      201   {object}  dto.AuditLogResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auditlogs/all [get]
func (h *Handler) GetAuditLogs(c *gin.Context) {
	page := int(utils.ConvertStringToInt(c.DefaultQuery("page", "1")))
	limit := int(utils.ConvertStringToInt(c.DefaultQuery("limit", "20")))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	resp, err := h.Deps.Service.GetAuditLog(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get audit logs"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
