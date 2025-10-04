package model

import (
	"database/sql"
	"time"
)

type UserAddress struct {
	ID              int          `gorm:"primaryKey"`
	UserID          int          `gorm:"type:int;not null"`
	ZipCode         string       `gorm:"type:varchar(16);not null"`
	Country         string       `gorm:"type:varchar(64);not null"`
	Province        string       `gorm:"type:varchar(64);not null"`
	City            string       `gorm:"type:varchar(64);not null"`
	Detail          string       `gorm:"type:varchar(255);not null"`
	FirstName       string       `gorm:"type:varchar(64);not null"`
	LastName        string       `gorm:"type:varchar(64);not null"`
	ContactPhone    string       `gorm:"type:varchar(32);not null"`
	DefaultMarkTime int64        `gorm:"column:default_mark_time;not null;default:0"`
	CreatedAt       time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       sql.NullTime `gorm:"column:deleted_at;"`
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
