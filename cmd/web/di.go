package main

import (
	"context"
	"errors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kevinrobayna/goprod/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var WebModule = fx.Module("web",
	fx.Provide(provideRouter),
	fx.Provide(provideWebHandler),
	fx.Invoke(invokeRoutes),
	fx.Invoke(invokeHttpServer),
)

func provideRouter(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return r
}

func provideWebHandler(logger *zap.Logger, router *gin.Engine, service internal.IService) IRoutes {
	return &WebHandler{
		Logger:  logger,
		R:       router,
		service: service,
	}
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
