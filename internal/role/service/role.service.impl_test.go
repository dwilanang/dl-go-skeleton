package service

import (
	"errors"
	"testing"

	"github.com/dwilanang/psp/internal/role/dto"
	"github.com/dwilanang/psp/internal/role/model"
	mockrepo "github.com/dwilanang/psp/internal/role/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	expected := []*model.Role{
		{ID: 1, Name: "Admin", Privilege: "all"},
		{ID: 2, Name: "User", Privilege: "read"},
	}

	mockRepo.EXPECT().Fetch().Return(expected, nil)

	resp, err := svc.GetAll()

	data := resp.Data.([]*model.Role)

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(data))
}

func TestService_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	mockRepo.EXPECT().Fetch().Return(nil, errors.New("fetch error"))

	_, err := svc.GetAll()
	assert.Error(t, err)
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.RoleRequest{Name: "Manager", Privilege: "manage", By: int64(1)}

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	err := svc.Create(req)
	assert.NoError(t, err)
}

func TestService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.RoleRequest{Name: "Manager", Privilege: "manage", By: int64(1)}

	mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("create error"))

	err := svc.Create(req)
	assert.Error(t, err)
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.RoleRequest{ID: 1, Name: "Updated", Privilege: "edit", By: int64(1)}

	mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	err := svc.Update(req)
	assert.NoError(t, err)
}

func TestService_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.RoleRequest{ID: 1, Name: "Updated", Privilege: "edit", By: int64(1)}

	mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))

	err := svc.Update(req)
	assert.Error(t, err)
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	mockRepo.EXPECT().Delete(int64(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestService_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	mockRepo.EXPECT().Delete(int64(1)).Return(errors.New("delete error"))

	err := svc.Delete(1)
	assert.Error(t, err)
}
