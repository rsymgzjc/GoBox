package main

import (
	"log"

	"gobox/backend/internal/config"
	"gobox/backend/internal/database"
	"gobox/backend/internal/routes"

	"go.uber.org/zap"
)

func main() {
	cfg, logger, err := config.Load() //加载配置文件
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	defer logger.Sync() //延迟执行

	db, err := database.Connect(cfg, logger) //链接数据库
	if err != nil {
		logger.Fatal("connect database", zap.Error(err))
	}

	router, err := routes.NewRouter(cfg, logger, db) //路由构建
	if err != nil {
		logger.Fatal("build router", zap.Error(err))
	}

	logger.Info("server starting", config.ZapString("addr", cfg.Server.Address()))
	if err := router.Run(cfg.Server.Address()); err != nil { //启动http服务器
		logger.Fatal("run server", zap.Error(err))
	}
}
