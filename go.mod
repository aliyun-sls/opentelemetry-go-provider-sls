module github.com/aliyun-sls/opentelemetry-go-provider-sls

go 1.15

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/sethvargo/go-envconfig v0.8.2
	github.com/tklauser/numcpus v0.5.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/host v0.35.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.35.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.31.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.10.0
	go.opentelemetry.io/otel/metric v0.31.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/sdk/export/metric v0.28.0
	go.opentelemetry.io/otel/sdk/metric v0.31.0
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220913175220-63ea55921009 // indirect
	google.golang.org/genproto v0.0.0-20220913154956-18f8339a66a5 // indirect
	google.golang.org/grpc v1.49.0
)
