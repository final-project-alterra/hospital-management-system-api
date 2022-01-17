package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		ENV.DB_USER,
		ENV.DB_PASSWORD,
		ENV.DB_HOST,
		ENV.DB_PORT,
		ENV.DB_NAME,
		ENV.DB_TIMEZONE,
	)

	dbLogger := logger.Default.LogMode(logger.Info)
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		panic(err)
	}

	DB = connection
}
