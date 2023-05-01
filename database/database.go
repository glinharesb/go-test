package database

import (
	"fmt"
	"go-test/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetConfig().DbUser,
		config.GetConfig().DbPass,
		config.GetConfig().DbHost,
		config.GetConfig().DbPort,
		config.GetConfig().DbDatabase)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func GetDb() *gorm.DB {
	return db.Debug()
}
