package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/keenJoe/go-url-shortener/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(conf *config.Config) error {
	// 打印配置信息（隐藏密码）
	log.Printf("正在连接数据库: host=%s, port=%d, user=%s, dbname=%s, password=%s",
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Username,
		conf.Database.DBName,
		conf.Database.Password)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName)

	// 尝试连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接测试失败: %v (DSN: %s)",
			err,
			strings.Replace(dsn, conf.Database.Password, "****", 1))
	}

	// 获取底层的sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库ping测试失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns) // 空闲连接数
	sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns) // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour)               // 连接最大生命周期

	DB = db
	return nil
}
