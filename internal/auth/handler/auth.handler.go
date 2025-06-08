package handler

import (
	"net/http"

	"github.com/dwilanang/psp/internal/auth/dto"
	"github.com/dwilanang/psp/internal/auth/service"
	utilrequest "github.com/dwilanang/psp/utils/request"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		Service: svc,
	}
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.AuthRequest  true  "Login credentials"
// @Success      200   {object}  dto.AuthResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var ar dto.AuthRequest
	if err := c.ShouldBindJSON(&ar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utilrequest.ValidateRequest(err)})
		return
	}

	resp, err := h.Service.Login(&ar)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
