package service

import (
	"fmt"

	"github.com/dwilanang/psp/internal/user/dto"
	"github.com/dwilanang/psp/internal/user/model"
	"github.com/dwilanang/psp/internal/user/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) *service {
	return &service{repo: r}
}

// Register implements the Service interface.
func (s *service) Register(request *dto.UserRequest) (dto.UserResponse, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, err
	}

	id := uuid.New()

	user := &model.User{
		UUID:         id.String(),
		Username:     request.Username,
		PasswordHash: string(hashed),
		FullName:     request.FullName,
		RoleID:       request.RoleID,
	}

	err = s.repo.Create(user)
	if err != nil {
		fmt.Println("s.repo.Create() error: ", err)
	}

	return dto.UserResponse{
		Data: dto.UserData{
			Username: user.Username,
			FullName: user.FullName,
		},
	}, err
}

// CreateSalary implements the Service interface.
func (s *service) CreateSalary(request *dto.UserSalaryRequest) (dto.UserSalaryResponse, error) {

	us := &model.UserSalary{
		UserID:        request.UserID,
		Amount:        request.Amount,
		EffectiveFrom: request.EffectiveFrom,
	}

	err := s.repo.CreateSalary(us)
	if err != nil {
		fmt.Println("s.repo.CreateSalary() error: ", err)
	}

	return dto.UserSalaryResponse{
		FullName: us.FullName,
		Amount:   us.Amount,
	}, err
}
