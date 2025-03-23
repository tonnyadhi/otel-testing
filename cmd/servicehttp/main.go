package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/api"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/config"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/database"
	"github.com/pintu-crypto/sre-playground/otel-testing/services/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	if cfg.Otel.Enable == true {
		otelTracer := telemetry.InitTracer(cfg.Otel.CollectorEndpoint, cfg.Otel.ServiceName, cfg.Otel.InsecureMode)
		otelMetric := telemetry.InitMetric(cfg.Otel.CollectorEndpoint, cfg.Otel.ServiceName, cfg.Otel.InsecureMode)
		log.Println("Tracer Initiated")
		defer otelTracer(context.Background())
		defer otelMetric(context.Background())
	}

	postgres, err := database.ConnectPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer postgres.DB.Close()

	queries := database.New(postgres.DB)
	httpService := api.NewService(queries)

	router := gin.Default()
	router.Use(otelgin.Middleware(cfg.Otel.ServiceName))
	router.Use(telemetry.MetricMiddleware(
		cfg.Otel.ServiceName,
		telemetry.WithAttributes(func(serverName, route string, request *http.Request) []attribute.KeyValue {
			return append(telemetry.DefaultMetricsAttributes(serverName, route, request), attribute.String("x-org", "pintu-crypto"))
		}),
	))
	router.Use(telemetry.LogrusMiddleware)

	httpService.RegisterHandlers(router)

	router.Run()
}
