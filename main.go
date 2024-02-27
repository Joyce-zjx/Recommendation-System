package main

import (
	"Recommendation-System/envconfig"
	"Recommendation-System/repository"
	"Recommendation-System/routes"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var env envconfig.Env
var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		logger.Error("failed to initialize logger ", zap.String("error message", err.Error()))
	}
	err = envconfig.Process(&env)
	if err != nil {
		logger.Error("failed to load config from env vars ", zap.String("error message", err.Error()))
	}
	gin.SetMode(gin.ReleaseMode)

	db, err := repository.NewGormDatabase(env.DATABASE_URL, env.Debug)
	if err != nil {
		logger.Error("failed to connect to database ", zap.String("error message", err.Error()))
	}
	if err := db.AutoMigrate(); err != nil {
		logger.Error("failed to migrate database ", zap.String("error message", err.Error()))
	}
	logger.Info("Finished migrating database")

	server := routes.Register(gin.Default(), logger, db, &env)
	go func() {
		server.Run(":5500")
	}()
	logger.Info("Server started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown servers...")
}
