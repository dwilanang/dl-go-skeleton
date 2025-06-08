package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auth/util"
	"github.com/dwilanang/psp/internal/role"
	"github.com/dwilanang/psp/internal/role/dto"
	"github.com/dwilanang/psp/utils"
	utilrequest "github.com/dwilanang/psp/utils/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Deps role.Dependencies
}

func NewHandler(deps role.Dependencies) *Handler {
	return &Handler{
		Deps: deps,
	}
}

// GetAll godoc
// @Security BearerAuth
// @Summary Get all roles
// @Description Get list of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {array} dto.RoleResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /roles/all [get]
func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.Deps.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch roles"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Create godoc
// @Security BearerAuth
// @Summary      Create roles
// @Description  Create a new roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RoleRequest  true  "Roles create payload"
// @Success      201   {object}  dto.RoleResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /roles/create [post]
func (h *Handler) Create(c *gin.Context) {
	var rr dto.RoleRequest
	if err := c.ShouldBindJSON(&rr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	by, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	rr.By = by

	if err := h.Deps.Service.Create(&rr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create role"})
		return
	}

	c.JSON(http.StatusCreated, rr)
}

// Update godoc
// @Security BearerAuth
// @Summary      Update roles
// @Description  Update a new roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Role ID"
// @Param        body  body      dto.RoleRequest  true  "Roles update payload"
// @Success      201   {object}  dto.RoleResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /roles/update/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var rr dto.RoleRequest
	if err := c.ShouldBindJSON(&rr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	by, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	rr.By = by

	idInt := utils.ConvertStringToInt(id)
	rr.ID = idInt

	if err := h.Deps.Service.Update(&rr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update role"})
		return
	}
	c.JSON(http.StatusOK, rr)
}

// Delete godoc
// @Security BearerAuth
// @Summary      Delete roles
// @Description  Delete roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id    path      int             true  "Role ID"
// @Success      201   {object}  dto.RoleResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /roles/delete/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	rr := dto.RoleResponse{}

	idInt := utils.ConvertStringToInt(id)

	if err := h.Deps.Service.Delete(idInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete role"})
		return
	}
	rr.Message = "Roles has been deleted."
	c.JSON(http.StatusOK, rr)
}
