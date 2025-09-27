package model

import "time"

type UserActivation struct {
	ID        int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID    int       `gorm:"type:int;not null;index"`
	Code      string    `gorm:"type:varchar(8);not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"type:datetime;not null"`
}

// TableName sets the insert table name for this struct type
func (UserActivation) TableName() string {
	return "user_activations"
}
