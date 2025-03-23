package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

type Option interface {
	apply(cfg *metricConfig)
}

type optionFunc func(cfg *metricConfig)

func (f optionFunc) apply(cfg *metricConfig) {
	f(cfg)
}

// WithAttributes sets a func using which what attributes to be recorded can be specified.
// By default the DefaultAttributes is used
func WithAttributes(attributes func(serverName, route string, request *http.Request) []attribute.KeyValue) Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.attributes = attributes
	})
}

// WithRecordInFlight determines whether to record In Flight Requests or not
// By default the recordInFlight is true
func WithRecordInFlightDisabled() Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.recordInFlight = false
	})
}

// WithRecordDuration determines whether to record Duration of Requests or not
// By default the recordDuration is true
func WithRecordDurationDisabled() Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.recordDuration = false
	})
}

// WithRecordSize determines whether to record Size of Requests and Responses or not
// By default the recordSize is true
func WithRecordSizeDisabled() Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.recordSize = false
	})
}

// WithGroupedStatus determines whether to group the response status codes or not. If true 2xx, 3xx will be stored
// By default the groupedStatus is true
func WithGroupedStatusDisabled() Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.groupedStatus = false
	})
}

// WithRecorder sets a recorder for recording requests
// By default the open telemetry recorder is used
func WithRecorder(recorder Recorder) Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.recorder = recorder
	})
}

// WithShouldRecordFunc sets a func using which whether a record should be recorded
// By default the all api calls are recorded
func WithShouldRecordFunc(shouldRecord func(serverName, route string, request *http.Request) bool) Option {
	return optionFunc(func(cfg *metricConfig) {
		cfg.shouldRecord = shouldRecord
	})
}
