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

package main

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"math/rand"
	"time"

	"github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	slsConfig, err := provider.NewConfig(provider.WithServiceName("payment"),
		provider.WithServiceVersion("v0.1.0"),
		provider.WithTraceExporterEndpoint("stdout"),
		provider.WithMetricExporterEndpoint("stdout"),
		provider.WithSLSConfig("test-project", "test-otel", "access-key-id", "access-key-secret"))
	// 如果初始化失败则panic，可以替换为其他错误处理方式
	if err != nil {
		panic(err)
	}
	if err := provider.Start(slsConfig); err != nil {
		panic(err)
	}
	defer provider.Shutdown(slsConfig)

	mockTrace()
	mockMetrics()
}

func mockMetrics() {
	// 附加的Label信息
	labels := []attribute.KeyValue{
		attribute.String("label1", "value1"),
	}

	meter := global.Meter("ex.com/basic")

	c, _ := meter.AsyncFloat64().Counter("randval")

	// 观测值，用于定期获取某个计量值，回调函数每个上报周期会被调用一次
	go mockObserveMetric(c, labels)

	temperature, _ := meter.SyncFloat64().Counter("temperature")
	interrupts, _ := meter.SyncInt64().Counter("interrupts")

	ctx := context.Background()

	for {
		temperature.Add(ctx, 100+10*rand.NormFloat64(), labels...)
		interrupts.Add(ctx, int64(rand.Intn(100)), labels...)

		time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	}
}

func mockObserveMetric(c asyncfloat64.Counter, labels []attribute.KeyValue) {
	timer := time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
		c.Observe(context.Background(), rand.Float64(), labels...)
	}
	timer.Stop()
}

func mockTrace() {

	tracer := otel.Tracer("ex.com/basic")

	ctx0 := context.Background()

	ctx1, finish1 := tracer.Start(ctx0, "foo")
	defer finish1.End()

	ctx2, finish2 := tracer.Start(ctx1, "bar")
	defer finish2.End()

	ctx3, finish3 := tracer.Start(ctx2, "baz")
	defer finish3.End()

	ctx := ctx3
	getSpan(ctx)
	addAttribute(ctx)
	addEvent(ctx)
	recordException(ctx)
	createChild(ctx, tracer)
}

// example of getting the current span
// 获取当前的Span
func getSpan(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	fmt.Printf("current span: %v\n", span)
}

// example of adding an attribute to a span
// 向Span中添加属性值
func addAttribute(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.KeyValue{
		Key:   "label-key-1",
		Value: attribute.StringValue("label-value-1")})
}

// example of adding an event to a span
// 向Span中添加事件
func addEvent(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("event1", trace.WithAttributes(
		attribute.String("event-attr1", "event-string1"),
		attribute.Int64("event-attr2", 10)))
}

// example of recording an exception
// 记录Span结果以及错误信息
func recordException(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(errors.New("exception has occurred"))
	span.SetStatus(codes.Error, "internal error")
}

// example of creating a child span
// 创建子Span
func createChild(ctx context.Context, tracer trace.Tracer) {
	// span := trace.SpanFromContext(ctx)
	_, childSpan := tracer.Start(ctx, "child")
	defer childSpan.End()
	fmt.Printf("child span: %v\n", childSpan)
}
