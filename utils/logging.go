package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StructuredLoggerHandlerFunc() gin.HandlerFunc {
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

		logrus.WithFields(logrus.Fields{
			"Method":    reqMethod,
			"Uri":       reqUri,
			"Status":    statusCode,
			"Latency":   latencyTime,
			"ClientIP":  clientIP,
			"UserAgent": userAgent,
			"LogTrace":  LogTrace,
		}).Info("Process request")
		ctx.Next()
	}
}

func InitLogger(loglevel string, verbose bool) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if verbose {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else if LogLevel, err := logrus.ParseLevel(loglevel); err == nil {
		if loglevelSlice := []int{2, 3, 4, 5}; slices.Contains(loglevelSlice, int(LogLevel)) {
			logrus.SetLevel(LogLevel)
		}
	}
	logrus.Debugf("Initializing logger with %v", loglevel)
}
