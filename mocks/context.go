package mocks

import (
	"bytes"
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
