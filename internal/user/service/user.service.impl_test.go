package service

import (
	"testing"

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
