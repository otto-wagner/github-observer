package mocks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func MockContext(param, requestBody string) (context *gin.Context, responseRecorder *httptest.ResponseRecorder) {
	responseRecorder = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(responseRecorder)

	if len(requestBody) > 0 {
		context.Request = httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBuffer([]byte(requestBody))))
	} else {
		context.Request = httptest.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewBuffer([]byte(requestBody))))
	}
	if len(param) > 0 {
		context.Params = gin.Params{gin.Param{Key: "id", Value: param}}
	}

	return
}

func PerformRequest(router *gin.Engine, method string, path string, body string) (w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	b := strings.NewReader(body)
	req := httptest.NewRequest(method, path, b)
	router.ServeHTTP(w, req)
	return
}

func PerformRequestAuthorisation(router *gin.Engine, method string, path string, body string, secret string) (w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	b := strings.NewReader(body)
	req := httptest.NewRequest(method, path, b)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(body))
	mac := hex.EncodeToString(hash.Sum(nil))
	req.Header.Set("X-Hub-Signature-256", "sha256="+mac)
	router.ServeHTTP(w, req)
	return
}
