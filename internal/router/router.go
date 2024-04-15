package router

import (
	"github-observer/internal/listener"
	"github-observer/internal/watcher"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
)

func InitializeRoutes(e *gin.Engine, w watcher.IWatcher, l listener.IListener) {
	//err = e.SetTrustedProxies(configuration.TrustedProxies)
	//if err != nil {
	//	return
	//}

	addCorsMiddleware(e)

	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	appExecutors := viper.GetStringSlice("app.executors")
	for _, executor := range appExecutors {
		if executor == "prometheus" {
			e.GET("/metrics", func(c *gin.Context) {
				handler := promhttp.Handler()
				handler.ServeHTTP(c.Writer, c.Request)
			})
			break
		}
	}

	el := e.Group("/listen")
	el.POST("/action", l.Action)
	el.POST("/pullrequest", l.PullRequest)
	el.POST("/pullrequest/review", l.PullRequestReview)

	if w != nil {
		ew := e.Group("/watch")
		ew.GET("/actions", w.Actions)
		ew.GET("/pullrequests", w.PullRequests)
	}
}

func addCorsMiddleware(engine *gin.Engine) {
	configCors := cors.DefaultConfig()
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Authorization"}
	configCors.AllowOrigins = []string{"*"}
	engine.Use(cors.New(configCors))
}
