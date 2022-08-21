package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinrobayna/goprod/internal"
	"go.uber.org/zap"
	"net/http"
)

type IRoutes interface {
	Hello(c *gin.Context)
	Ping(c *gin.Context)
}

type Handler struct {
	Logger  *zap.Logger
	R       *gin.Engine
	service internal.IService
}

func invokeRoutes(r *gin.Engine, h IRoutes) {
	r.GET("/", h.Hello)
	r.GET("/ping", h.Ping)
}

func (h *Handler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
}

func (h *Handler) Ping(c *gin.Context) {
	product, _ := h.service.GetProducts()
	c.JSON(http.StatusOK, product)
}
