package main

import (
	"context"
	"errors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var WebModule = fx.Module("web",
	fx.Provide(provideRouter, provideWebHandler),
	fx.Invoke(invokeRoutes, invokeHttpServer),
)

// provideRouter HTTP Stuff
func provideRouter(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return r
}

func provideWebHandler(logger *zap.Logger, router *gin.Engine, db *gorm.DB) *WebHandler {
	return &WebHandler{
		Logger: logger,
		R:      router,
		db:     db,
	}
}

func invokeRoutes(r *gin.Engine, h *WebHandler) {
	r.GET("/", h.Hello)
	r.GET("/ping", h.Ping)
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
