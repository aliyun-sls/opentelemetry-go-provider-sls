module github.com/aliyun-sls/opentelemetry-go-provider-sls

go 1.15

require (
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20231016141302-07b5767bb0ed // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/sethvargo/go-envconfig v0.9.0
	github.com/shirou/gopsutil/v3 v3.23.11 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.45.0
	go.opentelemetry.io/contrib/instrumentation/host v0.46.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1
	go.opentelemetry.io/contrib/instrumentation/runtime v0.46.1
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/sdk/metric v1.21.0
	golang.org/x/net v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/grpc v1.60.1
)
