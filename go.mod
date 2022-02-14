module github.com/aliyun-sls/opentelemetry-go-provider-sls

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/sethvargo/go-envconfig v0.3.2
	github.com/shirou/gopsutil/v3 v3.21.11 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.28.0
	go.opentelemetry.io/contrib/instrumentation/host v0.27.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.27.0
	go.opentelemetry.io/otel v1.4.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.4.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.4.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.27.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.0
	go.opentelemetry.io/otel/metric v0.27.0
	go.opentelemetry.io/otel/sdk v1.4.0
	go.opentelemetry.io/otel/sdk/export/metric v0.27.0
	go.opentelemetry.io/otel/sdk/metric v0.27.0
	go.opentelemetry.io/otel/trace v1.4.0
	golang.org/x/net v0.0.0-20211201190559-0a0e4e1bb54c // indirect
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211129164237-f09f9a12af12 // indirect
	google.golang.org/grpc v1.44.0
)
