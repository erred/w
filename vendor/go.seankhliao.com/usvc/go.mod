module go.seankhliao.com/usvc

go 1.16

require (
	github.com/prometheus/client_golang v1.7.1
	github.com/rs/zerolog v1.18.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.12.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.12.0
	go.opentelemetry.io/contrib/propagators v0.12.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	go.seankhliao.com/apis v0.0.0-20200925201609-7c5465abda54
	google.golang.org/grpc v1.32.0
)
