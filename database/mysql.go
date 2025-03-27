package database

import (
	"fmt"
	"github.com/ichtrojan/cook/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mysql *gorm.DB

func ConnectMySQL(credentials config.Mysql) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		credentials.User, credentials.Pass, credentials.Host, credentials.Port, credentials.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	Mysql = db

	return nil
}
