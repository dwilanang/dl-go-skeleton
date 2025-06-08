package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	mockrepo "github.com/dwilanang/psp/internal/user/repository/mocks"

	"github.com/dwilanang/psp/internal/user/dto"
	"github.com/dwilanang/psp/internal/user/model"
)

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.UserRequest{
		Username: "testuser",
		Password: "secret123",
		FullName: "Test User",
		RoleID:   1,
	}

	// Match what is expected to be created
	mockRepo.
		EXPECT().
		Create(gomock.AssignableToTypeOf(&model.User{})).
		DoAndReturn(func(user *model.User) error {
			assert.Equal(t, req.Username, user.Username)
			assert.Equal(t, req.FullName, user.FullName)
			assert.Equal(t, req.RoleID, user.RoleID)
			assert.NotEmpty(t, user.PasswordHash)
			assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)))
			return nil
		})

	resp, err := svc.Register(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Username, resp.Data.Username)
	assert.Equal(t, req.FullName, resp.Data.FullName)
}

func TestService_CreateSalary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.UserSalaryRequest{
		UserID:        1,
		Amount:        7500000,
		EffectiveFrom: time.Now().Format("2006-01-02 15:04:05"),
	}

	mockRepo.
		EXPECT().
		CreateSalary(gomock.AssignableToTypeOf(&model.UserSalary{})).
		DoAndReturn(func(us *model.UserSalary) error {
			us.FullName = "John Doe" // Simulasi isian yang dikembalikan dari DB layer
			return nil
		})

	resp, err := svc.CreateSalary(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Amount, resp.Amount)
	assert.Equal(t, "John Doe", resp.FullName)
}

func TestService_CreateSalary_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	req := &dto.UserSalaryRequest{
		UserID:        int64(1),
		Amount:        float64(5000000),
		EffectiveFrom: time.Now().Format("2006-01-02 15:04:05"),
	}

	mockRepo.
		EXPECT().
		CreateSalary(gomock.Any()).
		Return(errors.New("db error"))

	resp, err := svc.CreateSalary(req)
	assert.Error(t, err)
	assert.Equal(t, float64(5000000), resp.Amount)
	assert.Empty(t, resp.FullName)
}
