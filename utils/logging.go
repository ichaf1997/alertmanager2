package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		latencyTime := fmt.Sprintf("%dms", endTime.Sub(startTime).Milliseconds())
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		userAgent := ctx.Request.UserAgent()
		clientIP := ctx.ClientIP()
		_, fn, line, _ := runtime.Caller(7)
		scriptName := filepath.Base(fn)
		LogTrace := fmt.Sprintf("%s:%d", scriptName, line)

		logger.WithFields(logrus.Fields{
			"Method":    reqMethod,
			"Uri":       reqUri,
			"Status":    statusCode,
			"Latency":   latencyTime,
			"ClientIP":  clientIP,
			"UserAgent": userAgent,
			"LogTrace":  LogTrace,
		}).Info("Accept http request")
		ctx.Next()
	}
}

func DefaultLogger(loglevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	if LogLevel, err := logrus.ParseLevel(loglevel); err == nil {
		logger.SetLevel(LogLevel)
		if LogLevel == logrus.TraceLevel {
			logger.SetReportCaller(true)
		}
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger

}
