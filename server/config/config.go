package config

import (
	"gorm.io/gorm"
)

var (
	GVA_DB    *gorm.DB
	LocalPath = "uploads/file"
	DbType    = "mysql"
)
