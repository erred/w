module go.seankhliao.com/usvc

go 1.14

require (
	github.com/rs/zerolog v1.18.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc v0.11.0
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace v0.11.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.11.0
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.11.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	google.golang.org/grpc v1.31.0
)
