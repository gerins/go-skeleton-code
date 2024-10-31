package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go-skeleton-code/pkg/log"

	"github.com/gin-gonic/gin"
)

// SetLogRequest sets up the logging request in the Gin context.
func SetLogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the existing context from the request
		ctx := c.Request.Context()

		// Add the log request to the context
		newCtxWithLog := log.NewRequest().SaveToContext(ctx)

		// Replace the request with the new context
		c.Request = c.Request.WithContext(newCtxWithLog)

		// Proceed to the next handler
		c.Next()
	}
}

// SaveLogRequest handles logging of request and response data.
func SaveLogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get context
		ctx := c.Request.Context()

		// Capture request body
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(io.NopCloser(bytes.NewBuffer(reqBody))) // Re-read buffer for next handler

		// Process request in the main handler
		c.Next()

		// Capture response
		respBody, _ := io.ReadAll(c.Request.Response.Body)

		// Save log request
		extractRequestData(ctx, c, reqBody, respBody)
		log.Context(ctx).Save()
	}
}

func extractRequestData(ctx context.Context, c *gin.Context, req, resp []byte) {
	requestLog := log.Context(ctx) // Get log request from context

	// Populate log information
	requestLog.IP = c.ClientIP()
	requestLog.Method = c.Request.Method
	requestLog.URL = c.Request.Host + c.Request.URL.String()
	requestLog.ReqHeader, requestLog.RespHeader = getHeader(c)
	requestLog.StatusCode = c.Writer.Status()

	// Set request body based on HTTP method
	if requestLog.Method == http.MethodGet || requestLog.Method == http.MethodDelete {
		requestLog.ReqBody = c.Request.URL.Query()
	} else if requestLog.ReqBody == nil {
		if err := json.Unmarshal(req, &requestLog.ReqBody); err != nil {
			requestLog.ReqBody = string(req)
		}
	}

	// Set response body
	if requestLog.RespBody == nil {
		if err := json.Unmarshal(resp, &requestLog.RespBody); err != nil {
			requestLog.RespBody = string(resp)
		}
	}
}

// getHeader extracts headers from the request or response.
func getHeader(c *gin.Context) (map[string][]string, map[string][]string) {
	var (
		req  = make(map[string][]string)
		resp = make(map[string][]string)
	)

	for k, v := range c.Request.Header {
		req[k] = v
	}
	for k, v := range c.Writer.Header() {
		resp[k] = v
	}

	return req, resp
}
