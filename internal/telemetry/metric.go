package telemetry

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
)

func InitMetric(endpoint, serviceName string, insecure bool) func(context.Context) error {

	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure == true {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithDialOption(),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithReconnectionPeriod(time.Second*30),
	)

	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
			attribute.String("org.name", "pintu-crypto"),
		),
		resource.WithProcess(),
		resource.WithProcessCommandArgs(),
		resource.WithOSType(),
		resource.WithOS(),
		resource.WithTelemetrySDK(),
	)

	if err != nil {
		log.Printf("Could Not Get Resources : %v\n", err)
	}

	otel.SetMeterProvider(
		metric.NewMeterProvider(
			metric.WithResource(resources),
			metric.WithReader(
				metric.NewPeriodicReader(
					exporter,
					metric.WithInterval(time.Second*30),
				),
			),
		),
	)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}

	return exporter.Shutdown
}
