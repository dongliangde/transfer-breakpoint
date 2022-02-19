package initialize

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化Sqllite数据库
func GormSqllite() *gorm.DB {
	if db, err := gorm.Open(sqlite.Open("file_upload.db"), &gorm.Config{}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		return db
	}
}
