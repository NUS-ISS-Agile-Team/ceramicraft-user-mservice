package repository

import (
	"database/sql"
	"fmt"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/config"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/repository/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

type TxBeginner interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}

var _ TxBeginner = (*gorm.DB)(nil) // Compile-time interface check

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.MySQLConfig.UserName,
		config.Config.MySQLConfig.Password,
		config.Config.MySQLConfig.Host,
		config.Config.MySQLConfig.Port,
		config.Config.MySQLConfig.DBName,
	)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(
		&model.User{},
		&model.UserActivation{},
		&model.UserAddress{},
	)
	if err != nil {
		panic(err)
	}
}
