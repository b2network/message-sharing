package initiates

import (
	"bsquared.network/message-sharing-applications/internal/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(database config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.UserName, database.Password, database.Host, database.Port, database.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(database.LogLevel)),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
