package main

import (
	"context"
	"errors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"net/http"
	"time"
)

func provideLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func provideRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return r
}

func invokeRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello"})
	})
	r.GET("/ping", func(c *gin.Context) {
		var product Product
		db.First(&product, 1)
		c.JSON(200, product)
	})
}

func invokeHttpServer(lc fx.Lifecycle, ginEngine *gin.Engine, logger *zap.Logger) {
	logger.Info("Executing httpServer.")
	server := http.Server{
		Addr:    ":8080",
		Handler: ginEngine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running server", zap.String("addr", server.Addr))
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("failed to listen and serve from server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func provideGormDB(logger *zap.Logger) (*gorm.DB, error) {
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault() // configure gorm to use this zapgorm.Logger for callbacks
	db, err := gorm.Open(sqlite.Open("./db.sqlite"), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func invokeMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Product{})
	if err != nil {
		return
	}

	db.Create(&Product{Code: "D42", Price: 100})
}

func FxEvent(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}

func main() {
	app := fx.New(
		fx.Provide(provideLogger, provideRouter, provideGormDB),
		fx.Invoke(invokeRoutes, invokeMigrate, invokeHttpServer),
		fx.WithLogger(FxEvent),
	)

	app.Run()
}
