package service

import (
	"context"
	"testing"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/dao/mocks"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/model"
	"github.com/stretchr/testify/mock"
)

func TestUserServiceImpl_Create(t *testing.T) {
	mockDao := new(mocks.UserDao)
	userService := &UserServiceImpl{mockDao}

	tests := []struct {
		email    string
		password string
		mockFunc func()
		wantErr  bool
	}{
		{
			email:    "test@example.com",
			password: "password123",
			mockFunc: func() {
				mockDao.On("CreateUser", mock.Anything, mock.MatchedBy(
					func(*model.User) bool { return true }, // 只校验类型，不比较字段
				)).Return(1, nil)
			},
			wantErr: false,
		},
		// {
		// 	email:    "invalid-email",
		// 	password: "password123",
		// 	mockFunc: func() {
		// 		mockDao.On("CreateUser", mock.Anything, "invalid-email", "password123").Return(fmt.Errorf("invalid email format"))
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	email:    "test@example.com",
		// 	password: "short",
		// 	mockFunc: func() {
		// 		mockDao.On("CreateUser", mock.Anything, "test@example.com", "short").Return(fmt.Errorf("password too short"))
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			tt.mockFunc()
			err := userService.Create(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			mockDao.AssertExpectations(t)
		})
	}
}
