package web

import (
	"context"
	"errors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kevinrobayna/goprod/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
	"time"
)

var Module = fx.Module("web",
	internal.Module,
	fx.Provide(provideRouter),
	fx.Provide(provideListener, providePort),
	fx.Provide(provideWebHandler),
	fx.Invoke(invokeRoutes),
	fx.Invoke(invokeHttpServer),
)

func provideRouter(logger *zap.Logger) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return r
}

func provideListener(logger *zap.Logger, cfg internal.AppConfig) net.Listener {
	logger.Info("using config", zap.String("addr", cfg.WebConfig.Addr))
	listener, err := net.Listen("tcp", cfg.WebConfig.Addr)
	if err != nil {
		logger.Fatal("failed to listen and serve from server", zap.Error(err))
	}
	return listener
}

type Port string

func providePort(listener net.Listener) Port {
	addr := strings.TrimPrefix(listener.Addr().String(), "[::]:")
	return Port(addr)
}

func provideWebHandler(logger *zap.Logger, router *gin.Engine, service internal.IService) IRoutes {
	return &Handler{
		Logger:  logger,
		R:       router,
		service: service,
	}
}

func invokeHttpServer(lc fx.Lifecycle, ginEngine *gin.Engine, listener net.Listener, logger *zap.Logger) {
	server := http.Server{
		Handler: ginEngine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running server", zap.String("addr", listener.Addr().String()))
			go func() {
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("failed to serve from server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
