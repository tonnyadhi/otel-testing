# Build
.PHONY: build
build:
	@echo "Building local image..."
	@docker build -t otel-testing:local .

# Stack
.PHONY:	stop
stop:
	@docker compose -f docker-compose.yml down -v

.PHONY: stop-otel-local
stop-otel-local:
	@docker compose -f docker-compose-otel-local.yml down -v

.PHONY: stop-otel-datadog
stop-otel-datadog:
	@docker compose -f docker-compose-otel-datadog.yml down -v

.PHONY:	start
start:
	@docker compose -f docker-compose.yml down -v
	@docker compose -f docker-compose.yml up --build -d


.PHONY:	start-otel-local
start-otel-local:
	@docker compose -f docker-compose-otel-local.yml down -v
	@docker compose -f docker-compose-otel-local.yml up --build

.PHONY:	start-otel-datadog
start-otel-datadog:
	@docker compose -f docker-compose-otel-datadog.yml down -v
	@docker compose -f docker-compose-otel-datadog.yml up --build
