package router

import (
	"github-listener/internal/listener"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitializeRoutes(e *gin.Engine, l listener.IListener) {
	//err = e.SetTrustedProxies(configuration.TrustedProxies)
	//if err != nil {
	//	return
	//}

	addCorsMiddleware(e)

	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	el := e.Group("/listen")
	el.POST("/actions", l.Action)
	el.POST("/pullrequests", l.PullRequest)
}

func addCorsMiddleware(engine *gin.Engine) {
	configCors := cors.DefaultConfig()
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Authorization"}
	configCors.AllowOrigins = []string{"*"}
	engine.Use(cors.New(configCors))
}
