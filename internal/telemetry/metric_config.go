package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type metricConfig struct {
	recordInFlight bool
	recordSize     bool
	recordDuration bool
	groupedStatus  bool
	recorder       Recorder
	attributes     func(serverName, route string, request *http.Request) []attribute.KeyValue
	shouldRecord   func(serverName, route string, request *http.Request) bool
}

func defaultMetricConfig() *metricConfig {
	return &metricConfig{
		recordInFlight: true,
		recordSize:     true,
		recordDuration: true,
		groupedStatus:  true,
		recorder:       nil,
		attributes:     DefaultMetricsAttributes,
		shouldRecord: func(_, _ string, _ *http.Request) bool {
			return true
		},
	}
}

var DefaultMetricsAttributes = func(serverName, route string, request *http.Request) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPMethodKey.String(request.Method),
	}

	if serverName != "" {
		attrs = append(attrs, semconv.HTTPServerNameKey.String(serverName))
	}
	if route != "" {
		attrs = append(attrs, semconv.HTTPRouteKey.String(route))
	}
	return attrs
}
