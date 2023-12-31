package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"io"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		start := time.Now()
		c.Next()

		cost := time.Since(start)
		respStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", respStatus),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", helpers.MicrosecondsStr(cost)),
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			logFields = append(logFields, zap.String("Request Body", string(requestBody)))
			logFields = append(logFields, zap.String("Response Body", w.body.String()))
		}

		if respStatus > 400 && respStatus <= 499 {
			logger.Warn("HTTP Warning "+cast.ToString(respStatus), logFields...)
		} else if respStatus >= 500 && respStatus <= 599 {
			logger.Error("HTTP Error "+cast.ToString(respStatus), logFields...)
		} else {
			logger.Debug("HTTP Access Log", logFields...)
		}
	}
}
