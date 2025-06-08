package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auth/util"
	"github.com/dwilanang/psp/internal/user"
	"github.com/dwilanang/psp/internal/user/dto"
	utilrequest "github.com/dwilanang/psp/utils/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Deps user.Dependencies
}

func NewHandler(deps user.Dependencies) *Handler {
	return &Handler{
		Deps: deps,
	}
}

// Register godoc
// @Security BearerAuth
// @Summary      Register user
// @Description  Create a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      dto.UserRequest  true  "User registration payload"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var ur dto.UserRequest
	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	by, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	ur.By = by

	resp, err := h.Deps.Service.Register(&ur)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// CreateSalary godoc
// @Security BearerAuth
// @Summary      CreateSalary employee salary
// @Description  Create a new employee salary
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      dto.UserSalaryRequest  true  "Employee salary payload"
// @Success      201   {object}  dto.UserSalaryResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/salary [post]
func (h *Handler) CreateSalary(c *gin.Context) {
	var us dto.UserSalaryRequest
	if err := c.ShouldBindJSON(&us); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	by, err := util.GetClaimsID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims"})
		return
	}
	us.By = by

	resp, err := h.Deps.Service.CreateSalary(&us)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user salary"})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
