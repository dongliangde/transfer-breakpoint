package initialize

import (
	"log"
	"os"
	"transfer-breakpoint/config"
	"transfer-breakpoint/model"

	"gorm.io/gorm"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch config.DbType {
	case "mysql":
		return GormMysql()
	case "sqllite":
		return GormSqllite()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.FileUpload{},
	)
	if err != nil {
		log.Println("register table failed", err)
		os.Exit(0)
	}
	log.Println("register table success")
}
