package config

import (
	"fmt"
	"github.com/simabdi/gofiber-exception/exception"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUsername, DbPassword, DbHost, DbPort, DbDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"dsn": dsn,
			"db":  db,
		}).Info("Connection")

		panic(exception.Error(err))
	}

	DB = db
}
