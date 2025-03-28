FROM golang:1.22.3-alpine3.19 as builder
WORKDIR /otel

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o servicehttp ./cmd/servicehttp/main.go

FROM alpine:3.19
WORKDIR /bin
COPY --from=builder /otel/servicehttp /bin/servicehttp
ENV GIN_MODE=release
CMD /bin/servicehttp