package telemetry

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = &logrus.Logger{
	Out:   os.Stderr,
	Hooks: make(logrus.LevelHooks),
	Level: logrus.InfoLevel,
	Formatter: &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "log.level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function.name",
		},
		TimestampFormat: time.RFC3339Nano,
	},
}

func contextLogger(c *gin.Context) logrus.FieldLogger {
	return logger
}

func LogrusMiddleware(ctx *gin.Context) {
	start := time.Now()
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	if rawQuery := ctx.Request.URL.RawQuery; rawQuery != "" {
		path += "?" + rawQuery
	}
	ctx.Next()
	status := ctx.Writer.Status()
	contextLogger(ctx).Infof("%s %s %d %s", method, path, status, time.Since(start))
}
