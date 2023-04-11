// Copyright The AliyunSLS Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlpTraceGrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	slsProjectHeader         = "x-sls-otel-project"
	slsInstanceIDHeader      = "x-sls-otel-instance-id"
	slsAccessKeyIDHeader     = "x-sls-otel-ak-id"
	slsAccessKeySecretHeader = "x-sls-otel-ak-secret"
	slsSecurityTokenHeader   = "x-sls-otel-token"
)

// Option configures the sls otel provider
type Option func(*Config)

// WithMetricExporterEndpoint configures the endpoint for sending metrics via OTLP
// 配置Metric的输出地址，如果配置为空则禁用Metric功能，配置为stdout则打印到标准输出用于测试
func WithMetricExporterEndpoint(url string) Option {
	return func(c *Config) {
		c.MetricExporterEndpoint = url
	}
}

// WithTraceExporterEndpoint configures the endpoint for sending traces via OTLP
// 配置Trace的输出地址，如果配置为空则禁用Trace功能，配置为stdout则打印到标准输出用于测试
func WithTraceExporterEndpoint(url string) Option {
	return func(c *Config) {
		c.TraceExporterEndpoint = url
	}
}

// WithServiceName configures a "service.name" resource label
// 配置服务名称
func WithServiceName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

func WithServiceNamespace(namespace string) Option {
	return func(c *Config) {
		c.ServiceNamespace = namespace
	}
}

// WithServiceVersion configures a "service.version" resource label
// 配置版本号
func WithServiceVersion(version string) Option {
	return func(c *Config) {
		c.ServiceVersion = version
	}
}

// WithTraceExporterInsecure permits connecting to the trace endpoint without a certificate
// 配置是否禁用SSL，如果输出到SLS，则必须打开SLS
func WithTraceExporterInsecure(insecure bool) Option {
	return func(c *Config) {
		c.TraceExporterEndpointInsecure = insecure
	}
}

// WithMetricExporterInsecure permits connecting to the metric endpoint without a certificate
// 配置是否禁用SSL，如果输出到SLS，则必须打开SLS
func WithMetricExporterInsecure(insecure bool) Option {
	return func(c *Config) {
		c.MetricExporterEndpointInsecure = insecure
	}
}

// WithResourceAttributes configures attributes on the resource
// 配置上传附加的一些tag信息，例如环境、可用区等
func WithResourceAttributes(attributes map[string]string) Option {
	return func(c *Config) {
		c.resourceAttributes = attributes
	}
}

// WithResource configures attributes on the resource
// 配置上传附加的一些tag信息，例如环境、可用区等
func WithResource(resource *resource.Resource) Option {
	return func(c *Config) {
		c.Resource = resource
	}
}

// WithErrorHandler Configures a global error handler to be used throughout an OpenTelemetry instrumented project.
// See "go.opentelemetry.io/otel"
// 配置OpenTelemetry错误处理函数
func WithErrorHandler(handler otel.ErrorHandler) Option {
	return func(c *Config) {
		c.errorHandler = handler
	}
}

// WithMetricReportingPeriod configures the metric reporting period,
// how often the controller collects and exports metric data.
// 配置Metric导出间隔，默认为30s
func WithMetricReportingPeriod(p time.Duration) Option {
	return func(c *Config) {
		c.MetricReportingPeriod = fmt.Sprint(p)
	}
}

// WithSLSConfig configures sls project, instanceID, accessKeyID, accessKeySecret to send data to sls directly
// 配置输出到SLS的信息，包括 project, instanceID, accessKeyID, accessKeySecret
func WithSLSConfig(project, instanceID, accessKeyID, accessKeySecret string) Option {
	return func(c *Config) {
		c.Project, c.InstanceID, c.AccessKeyID, c.AccessKeySecret = project, instanceID, accessKeyID, accessKeySecret
	}
}

func WithIDGenerator(generator sdktrace.IDGenerator) Option {
	return func(config *Config) {
		if generator != nil {
			config.IDGenerator = generator
		}
	}
}

// Config configure for sls otel
type Config struct {
	TraceExporterEndpoint          string `env:"SLS_OTEL_TRACE_ENDPOINT,default=stdout"`
	TraceExporterEndpointInsecure  bool   `env:"SLS_OTEL_TRACE_INSECURE,default=false"`
	MetricExporterEndpoint         string `env:"SLS_OTEL_METRIC_ENDPOINT,default=stdout"`
	MetricExporterEndpointInsecure bool   `env:"SLS_OTEL_METRIC_INSECURE,default=false"`
	MetricReportingPeriod          string `env:"SLS_OTEL_METRIC_EXPORT_PERIOD,default=30s"`
	ServiceName                    string `env:"SLS_OTEL_SERVICE_NAME"`
	ServiceNamespace               string `env:"SLS_OTEL_SERVICE_NAMESPACE"`
	ServiceVersion                 string `env:"SLS_OTEL_SERVICE_VERSION,default=v0.1.0"`
	Project                        string `env:"SLS_OTEL_PROJECT"`
	InstanceID                     string `env:"SLS_OTEL_INSTANCE_ID"`
	AccessKeyID                    string `env:"SLS_OTEL_ACCESS_KEY_ID"`
	AccessKeySecret                string `env:"SLS_OTEL_ACCESS_KEY_SECRET"`
	AttributesEnvKeys              string `env:"SLS_OTEL_ATTRIBUTES_ENV_KEYS"`
	IDGenerator                    sdktrace.IDGenerator

	Resource *resource.Resource

	resourceAttributes map[string]string
	errorHandler       otel.ErrorHandler
	stop               []func()
}

func parseEnvKeys(c *Config) {
	if c.AttributesEnvKeys == "" {
		return
	}
	envKeys := strings.Split(c.AttributesEnvKeys, "|")
	for _, key := range envKeys {
		key = strings.TrimSpace(key)
		value := os.Getenv(key)
		if value != "" {
			c.resourceAttributes[key] = value
		}
	}
}

// 默认使用本机hostname作为hostname
func getDefaultResource(c *Config) *resource.Resource {
	hostname, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.ServiceName),
		semconv.HostNameKey.String(hostname),
		semconv.ServiceNamespaceKey.String(c.ServiceNamespace),
		semconv.ServiceVersionKey.String(c.ServiceVersion),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.ProcessCommandKey.String(os.Args[0]),
	)
}

func mergeResource(c *Config) error {
	var e error
	if c.Resource, e = resource.Merge(getDefaultResource(c), c.Resource); e != nil {
		return e
	}

	r := resource.Environment()
	if c.Resource, e = resource.Merge(c.Resource, r); e != nil {
		return e
	}

	var keyValues []attribute.KeyValue
	for key, value := range c.resourceAttributes {
		keyValues = append(keyValues, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: attribute.StringValue(value),
		})
	}
	newResource := resource.NewWithAttributes(semconv.SchemaURL, keyValues...)
	if c.Resource, e = resource.Merge(c.Resource, newResource); e != nil {
		return e
	}
	return nil
}

// 初始化Exporter，如果otlpEndpoint传入的值为 stdout，则默认把信息打印到标准输出用于调试
func (c *Config) initOtelExporter(otlpEndpoint string, insecure bool) (trace.SpanExporter, metric.Exporter, func(), error) {
	var traceExporter trace.SpanExporter
	var metricsExporter metric.Exporter
	var err error

	var exporterStop = func() {
		if traceExporter != nil {
			traceExporter.Shutdown(context.Background())
		}
	}

	if otlpEndpoint == "stdout" {
		// 使用Pretty的打印方式
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, nil, nil, err
		}
		enc := json.NewEncoder(os.Stdout)
		metricsExporter, err = stdoutmetric.New(stdoutmetric.WithEncoder(enc))
	} else if otlpEndpoint != "" {
		headers := map[string]string{}
		if c.Project != "" && c.InstanceID != "" {
			headers = map[string]string{
				slsProjectHeader:         c.Project,
				slsInstanceIDHeader:      c.InstanceID,
				slsAccessKeyIDHeader:     c.AccessKeyID,
				slsAccessKeySecretHeader: c.AccessKeySecret,
			}
		}

		// 使用GRPC方式导出数据
		traceSecureOption := otlpTraceGrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			traceSecureOption = otlpTraceGrpc.WithInsecure()
		}
		traceExporter, err = otlptrace.New(context.Background(),
			otlpTraceGrpc.NewClient(otlpTraceGrpc.WithEndpoint(otlpEndpoint),
				traceSecureOption,
				otlpTraceGrpc.WithHeaders(headers),
				otlpTraceGrpc.WithCompressor(gzip.Name)))
		if err != nil {
			return nil, nil, nil, err
		}

		metricSecureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			metricSecureOption = otlpmetricgrpc.WithInsecure()
		}

		metricsExporter, err = otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithEndpoint(otlpEndpoint),
			metricSecureOption, otlpmetricgrpc.WithHeaders(headers), otlpmetricgrpc.WithCompressor(gzip.Name))
	}

	return traceExporter, metricsExporter, exporterStop, nil
}

// 初始化Metrics，默认30秒导出一次Metrics
// 默认该函数导出主机和Golang runtime基础指标
func (c *Config) initMetric(metricsExporter metric.Exporter, stop func()) error {
	if metricsExporter == nil {
		return nil
	}
	period, err := time.ParseDuration(c.MetricReportingPeriod)
	if err != nil {
		period = time.Second * 30
	}

	reader := metric.NewPeriodicReader(metricsExporter, metric.WithInterval(period))

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(c.Resource))
	global.SetMeterProvider(meterProvider)

	// 默认集成主机基础指标
	if err := host.Start(host.WithMeterProvider(meterProvider)); err != nil {
		return err
	}
	// 默认集成Golang runtime指标
	err = runtime.Start(runtime.WithMeterProvider(meterProvider), runtime.WithMinimumReadMemStatsInterval(time.Second))
	c.stop = append(c.stop, func() {
		meterProvider.Shutdown(context.Background())
		stop()
	})
	return err
}

// 初始化Traces，默认全量上传
func (c *Config) initTracer(traceExporter trace.SpanExporter, stop func(), config *Config) error {
	if traceExporter == nil {
		return nil
	}
	// 建议使用AlwaysSample全量上传Trace数据，若您的数据太多，可以使用sdktrace.ProbabilitySampler进行采样上传
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			traceExporter,
		),
		sdktrace.WithIDGenerator(config.IDGenerator),
		sdktrace.WithResource(c.Resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	c.stop = append(c.stop, func() {
		tp.Shutdown(context.Background())
		stop()
	})
	return nil
}

// IsValid check config and return error if config invalid
func (c *Config) IsValid() error {
	if c.ServiceName == "" {
		return errors.New("empty service name")
	}
	if c.ServiceVersion == "" {
		return errors.New("empty service version")
	}
	if (strings.Contains(c.TraceExporterEndpoint, "log.aliyuncs.com") && c.TraceExporterEndpointInsecure) ||
		(strings.Contains(c.MetricExporterEndpoint, "log.aliyuncs.com") && c.MetricExporterEndpointInsecure) {
		return errors.New("insecure grpc is not allowed when send data to sls directly")
	}
	if strings.Contains(c.TraceExporterEndpoint, "log.aliyuncs.com") || strings.Contains(c.MetricExporterEndpoint, "log.aliyuncs.com") {
		if c.Project == "" || c.InstanceID == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" {
			return errors.New("empty project, instanceID, accessKeyID or accessKeySecret when send data to sls directly")
		}
		if strings.ContainsAny(c.Project, "${}") ||
			strings.ContainsAny(c.InstanceID, "${}") ||
			strings.ContainsAny(c.AccessKeyID, "${}") ||
			strings.ContainsAny(c.AccessKeySecret, "${}") {
			return errors.New("invalid project, instanceID, accessKeyID or accessKeySecret when send data to sls directly, you should replace these parameters with actual values")
		}
	}
	return nil
}

// NewConfig create a config
func NewConfig(opts ...Option) (*Config, error) {
	var c Config

	// 1. load env config
	envError := envconfig.Process(context.Background(), &c)
	if envError != nil {
		return nil, envError
	}

	// 2. load code config
	for _, opt := range opts {
		opt(&c)
	}

	// 3. merge resource
	parseEnvKeys(&c)
	mergeResource(&c)
	return &c, c.IsValid()
}

// Start 初始化OpenTelemetry SDK，需要把 ${endpoint} 替换为实际的地址
// 如果填写为stdout则为调试默认，数据将打印到标准输出
func Start(c *Config) error {
	if c.errorHandler != nil {
		otel.SetErrorHandler(c.errorHandler)
	}
	traceExporter, _, traceExpStop, err := c.initOtelExporter(c.TraceExporterEndpoint, c.TraceExporterEndpointInsecure)
	if err != nil {
		return err
	}
	_, metricExporter, metricExpStop, err := c.initOtelExporter(c.MetricExporterEndpoint, c.MetricExporterEndpointInsecure)
	if err != nil {
		return err
	}
	err = c.initTracer(traceExporter, traceExpStop, c)
	if err != nil {
		return err
	}
	err = c.initMetric(metricExporter, metricExpStop)
	return err
}

// Shutdown 优雅关闭，将OpenTelemetry SDK内存中的数据发送到服务端
func Shutdown(c *Config) {
	for _, stop := range c.stop {
		stop()
	}
}
