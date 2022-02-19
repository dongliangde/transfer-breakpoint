package initialize

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化Mysql数据库
func GormMysql() *gorm.DB {
	mysqlConfig := mysql.Config{
		DSN:                       Dsn(), // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		return db
	}
}

func Dsn() string {
	dsn := "root:123456@tcp(127.0.0.1:3306)/file_upload?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}
