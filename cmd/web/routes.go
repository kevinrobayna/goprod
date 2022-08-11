package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goprod/internal/models"
	"gorm.io/gorm"
	"net/http"
)

type IRoutes interface {
	Hello(c *gin.Context)
	Ping(c *gin.Context)
}

type WebHandler struct {
	Logger *zap.Logger
	R      *gin.Engine
	db     *gorm.DB
}

func invokeRoutes(r *gin.Engine, h IRoutes) {
	r.GET("/", h.Hello)
	r.GET("/ping", h.Ping)
}

func (h *WebHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello"})
}

func (h *WebHandler) Ping(c *gin.Context) {
	var product models.Product
	h.db.First(&product, 1)
	c.JSON(200, product)
}
