package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github-observer/internal/listener"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func InitializeRoutes(e *gin.Engine, l listener.IListener, config Config) {
	err := e.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		return
	}

	addCorsMiddleware(e)

	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
	})

	e.GET("/configuration", func(c *gin.Context) {
		config.App.Listener.HmacSecret = config.App.Listener.HmacSecret[:4] + strings.Repeat("*", len(config.App.Listener.HmacSecret)-4)
		c.JSON(http.StatusOK, gin.H{"configuration": config})
	})

	var loggingExecutor bool
	for _, executor := range config.App.Executors {
		if executor == "prometheus" {
			e.GET("/metrics", gin.WrapH(promhttp.Handler()))
			break
		}
		if executor == "logger" {
			loggingExecutor = true
		}
	}

	logsGroup := e.Group("/logs")
	logsGroup.GET("/observer", func(c *gin.Context) {
		data, err := ReadLogData(ObserverFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, data)
	})
	for _, l := range config.App.Logger {
		switch l {
		case "listener":
			logsGroup.GET("/listener", func(c *gin.Context) {
				data, err := ReadLogData(ListenerFile)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.String(http.StatusOK, data)
			})
		case "executor":
			if loggingExecutor {
				logsGroup.GET("/executor", func(c *gin.Context) {
					data, err := ReadLogData(ExecutorFile)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.String(http.StatusOK, data)
				})
			}
		case "watcher":
			logsGroup.GET("/watcher", func(c *gin.Context) {
				data, err := ReadLogData(WatcherFile)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.String(http.StatusOK, data)
			})
		}
	}

	if l != nil {
		el := e.Group("/listen")
		el.Use(hmacMiddleware([]byte(config.App.Listener.HmacSecret)))
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

func hmacMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		macHeader := c.GetHeader("X-Hub-Signature-256")
		if len(macHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Hub-Signature-256"})
			return
		}
		actualMAC, err := hex.DecodeString(strings.Split(macHeader, "=")[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Hub-Signature-256"})
			return
		}
		mac := hmac.New(sha256.New, secretKey)
		mac.Write(body)
		expectedMAC := mac.Sum(nil)

		if !hmac.Equal(actualMAC, expectedMAC) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()
	}
}
