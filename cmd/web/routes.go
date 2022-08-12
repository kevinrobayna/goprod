package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goprod/internal/products"
	"net/http"
)

func NewRoutes(logger *zap.Logger, router *gin.Engine, service products.IService) IRoutes {
	return &WebHandler{
		Logger:  logger,
		R:       router,
		service: service,
	}
}

type IRoutes interface {
	Hello(c *gin.Context)
	Ping(c *gin.Context)
}

type WebHandler struct {
	Logger  *zap.Logger
	R       *gin.Engine
	service products.IService
}

func invokeRoutes(r *gin.Engine, h IRoutes) {
	r.GET("/", h.Hello)
	r.GET("/ping", h.Ping)
}

func (h *WebHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
}

func (h *WebHandler) Ping(c *gin.Context) {
	product, _ := h.service.GetProducts()
	c.JSON(http.StatusOK, product)
}
