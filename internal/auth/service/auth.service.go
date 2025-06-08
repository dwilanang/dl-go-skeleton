package service

import (
	"github.com/dwilanang/psp/internal/auth/dto"
	"github.com/dwilanang/psp/utils/response"
)

type Service interface {
	Login(request *dto.AuthRequest) (response.ApiResponse, error)
}
