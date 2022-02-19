package main

import (
	"net/http"
	"time"
	"transfer-breakpoint/config"
	"transfer-breakpoint/initialize"
	"transfer-breakpoint/routers"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title 断点续传服务
// @version 1.0
// @description swagger Api
// @BasePath /
func main() {
	config.GVA_DB = initialize.Gorm()
	if config.GVA_DB != nil {
		initialize.RegisterTables(config.GVA_DB)
		// 程序结束前关闭数据库链接
		db, _ := config.GVA_DB.DB()
		defer db.Close()
	}

	s := &http.Server{
		Addr:           ":8080",
		Handler:        routers.InitRouter(),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
