package service

import (
	"context"
	"sync"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/dao"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/model"
)

type UserService interface {
	Create(ctx context.Context, email, password string) error
}

type UserServiceImpl struct {
	userDao dao.UserDao
}

var (
	registerOnce    = sync.Once{}
	registerService *UserServiceImpl
)

func GetRegisterService() *UserServiceImpl {
	registerOnce.Do(func() {
		if registerService == nil {
			registerService = &UserServiceImpl{
				userDao: dao.GetUserDao(),
			}
		}
	})
	return registerService
}

func (s *UserServiceImpl) Create(ctx context.Context, email, password string) error {
	_, err := s.userDao.CreateUser(ctx, &model.User{
		Email:    email,
		Password: password,
	})
	return err
}
