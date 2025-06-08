package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dwilanang/psp/config"
	"github.com/dwilanang/psp/internal/auth/dto"
	userrepository "github.com/dwilanang/psp/internal/user/repository"
	"github.com/dwilanang/psp/utils"
	"github.com/dwilanang/psp/utils/response"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	cfg      *config.Config
	userRepo userrepository.Repository
}

func NewService(cfg *config.Config, repo userrepository.Repository) Service {
	return &service{
		cfg:      cfg,
		userRepo: repo,
	}
}

func (s *service) Login(request *dto.AuthRequest) (response.ApiResponse, error) {
	u, err := s.userRepo.FindByUsername(request.Username)
	if err != nil {
		if err.Error() == "failed" {
			return response.ApiResponse{}, errors.New("Login failed")
		}
		return response.ApiResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(request.Password))
	if err != nil {
		fmt.Println("Login: CompareHashAndPassword ", err)
		return response.ApiResponse{}, errors.New("Login failed")
	}

	jwtExpiration := utils.ConvertStringToInt(s.cfg.JWTExpiration)

	claims := jwt.MapClaims{
		"uid":  u.ID,
		"sub":  u.UUID,
		"role": u.Role,
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpiration) * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := s.cfg.JWTSecret
	if secret == "" {
		secret = "secret"
	}

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return response.ApiResponse{}, err
	}

	return response.ApiResponse{
		Status:  true,
		Message: "Login succefully",
		Data: dto.AuthResponse{
			Type:   s.cfg.JWTType,
			Token:  tokenStr,
			Expire: fmt.Sprintf("%d Hour", jwtExpiration),
		},
	}, nil // or an error if authentication fails
}
