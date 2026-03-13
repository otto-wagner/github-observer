package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if path == "/health" {
			return
		}
		if gin.Mode() == gin.ReleaseMode && path == "/metrics" {
			return
		}

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		slog.Info("HTTP request",
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency,
			"ip", clientIP,
		)
	}
}

func Hmac(secretKey []byte) gin.HandlerFunc {
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
