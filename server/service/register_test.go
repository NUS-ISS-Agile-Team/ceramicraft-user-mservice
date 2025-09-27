package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	proxy_mock "github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/proxy/mocks"
	dao_mock "github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/dao/mocks"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	initEnv()
	ctx := context.Background()

	t.Run("Successful registration", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		userDao := new(dao_mock.UserDao)
		emailSender := new(proxy_mock.EmailService)
		service := &RegisterImpl{
			userDao:        userDao,
			userActivation: userActivationDao,
			emailService:   emailSender,
		}
		email := "test@example.com"
		password := "password123"
		userId := 1
		userDao.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
		userDao.On("CreateUser", mock.Anything, mock.MatchedBy(func(arg *model.User) bool {
			arg.ID = userId                                       // Simulate DB assigning ID
			return arg.Email == email && arg.Password != password // Password should be hashed
		})).Return(1, nil)
		userActivationDao.On("Replace", mock.Anything, mock.MatchedBy(func(arg *model.UserActivation) bool {
			return arg.UserID == userId && len(arg.Code) == 6
		})).Return(nil)
		emailSender.On("Send", mock.Anything, email, mock.Anything).Return(nil)
		err := service.Register(ctx, email, password)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		userDao.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
		userActivationDao.AssertCalled(t, "Replace", mock.Anything, mock.Anything)
		emailSender.AssertCalled(t, "Send", mock.Anything, email, mock.Anything)
		userDao.AssertExpectations(t)
		userActivationDao.AssertExpectations(t)
		emailSender.AssertExpectations(t)
	})

	t.Run("User already exists", func(t *testing.T) {
		userDao := new(dao_mock.UserDao)
		service := &RegisterImpl{
			userDao: userDao,
		}
		email := "test@example.com"
		password := "password123"
		userDao.On("GetUserByEmail", mock.Anything, email).Return(&model.User{Email: email, Status: model.UserStatusActive}, nil)
		err := service.Register(ctx, email, password)
		if err == nil || err.Error() != "user already exists" {
			t.Fatalf("Expected 'user already exists' error, got %v", err)
		}
	})

	t.Run("Database error on GetUserByEmail", func(t *testing.T) {
		userDao := new(dao_mock.UserDao)
		service := &RegisterImpl{
			userDao: userDao,
		}
		email := "test@example.com"
		password := "password123"
		userDao.On("GetUserByEmail", mock.Anything, email).Return(nil, assert.AnError)
		err := service.Register(ctx, email, password)
		if err == nil || !errors.Is(err, assert.AnError) {
			t.Fatalf("Expected database error, got %v", err)
		}
	})

	t.Run("Database error on CreateUser", func(t *testing.T) {
		userDao := new(dao_mock.UserDao)
		service := &RegisterImpl{
			userDao: userDao,
		}
		email := "test@example.com"
		userDao.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
		userDao.On("CreateUser", mock.Anything, mock.Anything).Return(-1, assert.AnError)
		err := service.Register(ctx, email, "password123")
		if err == nil || !errors.Is(err, assert.AnError) {
			t.Fatalf("Expected database error on CreateUser, got %v", err)
		}
	})

	t.Run("Database error on Replace activation", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		userDao := new(dao_mock.UserDao)
		service := &RegisterImpl{
			userDao:        userDao,
			userActivation: userActivationDao,
		}
		email := "test@example.com"
		userDao.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
		userDao.On("CreateUser", mock.Anything, mock.Anything).Return(1, nil)
		userActivationDao.On("Replace", mock.Anything, mock.Anything).Return(assert.AnError)
		err := service.Register(ctx, email, "password123")
		if err == nil || !errors.Is(err, assert.AnError) {
			t.Fatalf("Expected database error on Replace, got %v", err)
		}
	})

	t.Run("Email sending failure", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		userDao := new(dao_mock.UserDao)
		emailSender := new(proxy_mock.EmailService)
		service := &RegisterImpl{
			userDao:        userDao,
			userActivation: userActivationDao,
			emailService:   emailSender,
		}
		email := "test@example.com"
		userDao.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
		userDao.On("CreateUser", mock.Anything, mock.Anything).Return(1, nil)
		userActivationDao.On("Replace", mock.Anything, mock.Anything).Return(nil)
		emailSender.On("Send", mock.Anything, email, mock.Anything).Return(assert.AnError)
		err := service.Register(ctx, email, "password123")
		if err == nil || !errors.Is(err, assert.AnError) {
			t.Fatalf("Expected email sending error, got %v", err)
		}
	})
}

type fakeTx struct{ *gorm.DB }

func (f *fakeTx) Transaction(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return fn(f.DB) // Pass through the same DB instance
}

func initMemDb(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&model.User{}, &model.UserActivation{}))
	return db
}

func TestVerifyAndActivate(t *testing.T) {
	initEnv()
	ctx := context.Background()

	t.Run("Activation code expired", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		service := &RegisterImpl{
			userActivation: userActivationDao,
		}
		code := "expired-code"
		userActivationDao.On("GetByCode", mock.Anything, code).Return(&model.UserActivation{
			Code:      code,
			ExpiresAt: time.Now().Add(-time.Minute),
			UserID:    1,
		}, nil)
		err := service.VerifyAndActivate(ctx, code)
		if err == nil || err.Error() != "invalid or expired activation code" {
			t.Fatalf("Expected 'invalid or expired activation code' error, got %v", err)
		}
		userActivationDao.AssertCalled(t, "GetByCode", mock.Anything, code)
		userActivationDao.AssertExpectations(t)
	})

	t.Run("Successful activation", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		userDao := new(dao_mock.UserDao)
		emailSender := new(proxy_mock.EmailService)
		service := &RegisterImpl{
			userDao:        userDao,
			userActivation: userActivationDao,
			emailService:   emailSender,
			txBeginner:     &fakeTx{DB: initMemDb(t)},
		}
		validCode := "valid-code"
		userActivationDao.On("GetByCode", mock.Anything, mock.Anything).Return(&model.UserActivation{
			Code:      validCode,
			UserID:    1,
			ExpiresAt: time.Now().Add(time.Minute * 10),
		}, nil)

		userDao.On("UpdateUser", mock.Anything, mock.MatchedBy(func(arg *model.User) bool {
			return arg.Status == model.UserStatusActive
		}), mock.Anything).Return(nil)
		userActivationDao.On("DeleteByUserId", mock.Anything, 1, mock.Anything).Return(nil)
		err := service.VerifyAndActivate(ctx, validCode)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		userDao.AssertCalled(t, "UpdateUser", mock.Anything, mock.Anything, mock.Anything)
		userActivationDao.AssertCalled(t, "DeleteByUserId", mock.Anything, 1, mock.Anything)
		userDao.AssertExpectations(t)
		userActivationDao.AssertExpectations(t)
	})

	t.Run("Invalid activation code", func(t *testing.T) {
		userActivationDao := new(dao_mock.UserActivationDao)
		userActivationDao.On("GetByCode", mock.Anything, "invalid-code").Return(nil, nil)
		service := &RegisterImpl{
			userActivation: userActivationDao,
		}

		err := service.VerifyAndActivate(ctx, "invalid-code")
		if err == nil || err.Error() != "invalid or expired activation code" {
			t.Fatalf("Expected 'invalid or expired activation code' error, got %v", err)
		}
	})
}
func TestGetRegisterService(t *testing.T) {
	initEnv()
	t.Run("Singleton instance", func(t *testing.T) {
		service1 := GetRegisterService()
		service2 := GetRegisterService()
		if service1 != service2 {
			t.Fatalf("Expected the same instance, got different instances")
		}
	})

	t.Run("Service initialization", func(t *testing.T) {
		service := GetRegisterService()
		if service.userDao == nil {
			t.Fatal("Expected userDao to be initialized")
		}
		if service.userActivation == nil {
			t.Fatal("Expected userActivation to be initialized")
		}
		if service.emailService == nil {
			t.Fatal("Expected emailService to be initialized")
		}
		if service.txBeginner == nil {
			t.Fatal("Expected txBeginner to be initialized")
		}
	})
}
