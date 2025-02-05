package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Uuid           string `gorm:"type:varchar(100);unique"`
	FullName       string `gorm:"type:varchar(100)"`
	Email          string `gorm:"type:varchar(45);unique"`
	PhoneNumber    string `gorm:"type:varchar(13)"`
	Password       string `gorm:"type:varchar(200)"`
	ProfilePicture string `gorm:"type:varchar(100)"`
	Role           string `gorm:"type:varchar(15)"`
	Status         string `gorm:"type:varchar(25)"`
	RefID          uint
	RefType        string         `gorm:"type:varchar(35)"`
	LinkExpire     string         `gorm:"type:datetime"`
	CreatedAt      time.Time      `gorm:"<-:create;type:datetime(0)"`
	UpdatedAt      time.Time      `gorm:"<-:update;type:datetime(0)"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
