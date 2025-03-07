package database

import (
	"fmt"

	"github.com/keenJoe/go-url-shortener/config"
	"github.com/keenJoe/go-url-shortener/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(conf *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db

	// 自动迁移数据库表
	err = autoMigrate(db)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	return nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		// 这里可以添加其他模型
	)
}
