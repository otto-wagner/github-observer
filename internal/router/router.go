package router

import (
	"github-listener/internal/listener"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(e *gin.Engine, _ listener.IListener) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
