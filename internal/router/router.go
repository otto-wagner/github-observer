package router

import (
	"github-observer/conf"
	"github-observer/internal/listener"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func InitializeRoutes(e *gin.Engine, l listener.IListener, watcherEnabled bool) {
	//err = e.SetTrustedProxies(configuration.TrustedProxies)
	//if err != nil {
	//	return
	//}

	addCorsMiddleware(e)

	e.GET("/health", func(c *gin.Context) {
		zap.S().Info("Health check")
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	appExecutors := viper.GetStringSlice("app.executors")
	for _, executor := range appExecutors {
		if executor == "logging" {
			e.GET("/logs/executor", func(c *gin.Context) {
				data, err := ReadLogData(conf.ExecutorFile)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.String(http.StatusOK, data)
			})
		}
		if executor == "prometheus" {
			e.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}
	}

	if watcherEnabled {
		e.GET("/logs/watcher", func(c *gin.Context) {
			data, err := ReadLogData(conf.WatcherFile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.String(http.StatusOK, data)
		})
	}

	if l != nil {
		el := e.Group("/listen")
		el.POST("/workflow", l.Workflow)
		el.POST("/pullrequest", l.PullRequest)
		el.POST("/pullrequest/review", l.PullRequestReview)
	}

}

func addCorsMiddleware(engine *gin.Engine) {
	configCors := cors.DefaultConfig()
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Authorization"}
	configCors.AllowOrigins = []string{"*"}
	engine.Use(cors.New(configCors))
}

func ReadLogData(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading log file: %v", err)
		return "", err
	}
	return string(data), nil
}
