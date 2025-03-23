package telemetry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func MetricMiddleware(serviceName string, options ...Option) gin.HandlerFunc {
	metricCfg := defaultMetricConfig()
	for _, option := range options {
		option.apply(metricCfg)
	}

	recorder := metricCfg.recorder
	if recorder == nil {
		recorder = GetRecorder("")
	}

	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()

		route := ginCtx.FullPath()
		if len(route) <= 0 {
			route = "nonconfigured"
		}

		if !metricCfg.shouldRecord(serviceName, route, ginCtx.Request) {
			ginCtx.Next()
			return
		}

		start := time.Now()
		reqAttributes := metricCfg.attributes(serviceName, route, ginCtx.Request)

		if metricCfg.recordInFlight {
			recorder.AddInflightRequests(ctx, 1, reqAttributes)
			defer recorder.AddInflightRequests(ctx, -1, reqAttributes)
		}

		defer func() {
			resAttributes := append(reqAttributes[0:0], reqAttributes...)

			if metricCfg.groupedStatus {
				code := int(ginCtx.Writer.Status()/100) * 100
				resAttributes = append(resAttributes, semconv.HTTPStatusCodeKey.Int(code))
			} else {
				resAttributes = append(resAttributes, semconv.HTTPAttributesFromHTTPStatusCode(ginCtx.Writer.Status())...)
			}

			recorder.AddRequests(ctx, 1, resAttributes)

			if metricCfg.recordSize {
				requestSize := computeApproximateRequestSize(ginCtx.Request)
				recorder.ObserveHTTPRequestSize(ctx, requestSize, resAttributes)
				recorder.ObserveHTTPResponseSize(ctx, int64(ginCtx.Writer.Size()), resAttributes)
			}

			if metricCfg.recordDuration {
				recorder.ObserveHTTPRequestDuration(ctx, time.Since(start), resAttributes)
			}
		}()

		ginCtx.Next()
	}

}

func computeApproximateRequestSize(r *http.Request) int64 {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	//r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return int64(s)
}
