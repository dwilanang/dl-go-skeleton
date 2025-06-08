package service

import (
	"fmt"

	"github.com/dwilanang/psp/internal/role/dto"
	"github.com/dwilanang/psp/internal/role/model"
	"github.com/dwilanang/psp/internal/role/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) *service {
	return &service{repo: r}
}

// GetAll implements the Service interface.
func (s *service) GetAll() (dto.RoleResponse, error) {
	result, err := s.repo.Fetch()
	if err != nil {
		fmt.Println("s.repo.GetAll() error: ", err)
		return dto.RoleResponse{}, err
	}

	return dto.RoleResponse{
		Data: result,
	}, err
}

// Create implements the Service interface.
func (s *service) Create(request *dto.RoleRequest) error {

	role := &model.Role{
		Name:      request.Name,
		Privilege: request.Privilege,
		CreatedBy: request.By,
	}

	err := s.repo.Create(role)
	if err != nil {
		fmt.Println("s.repo.Create() error: ", err)
	}

	return err
}

// Update implements the Service interface.
func (s *service) Update(request *dto.RoleRequest) error {
	role := &model.Role{
		ID:        request.ID,
		Name:      request.Name,
		Privilege: request.Privilege,
		UpdatedBy: request.By,
	}

	err := s.repo.Update(role)
	if err != nil {
		fmt.Println("s.repo.Update() error: ", err)
	}

	return err
}

// Delete implements the Service interface.
func (s *service) Delete(id int64) error {

	err := s.repo.Delete(id)
	if err != nil {
		fmt.Println("s.repo.Delete() error: ", err)
	}

	return err
}
