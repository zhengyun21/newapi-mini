package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// TODO: 替换为你的 MySQL DSN，格式：user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "<your-mysql-dsn>"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库出现异常：" + err.Error())
	}
	db.AutoMigrate(&Channel{}, &Token{})

	DB = db
}
