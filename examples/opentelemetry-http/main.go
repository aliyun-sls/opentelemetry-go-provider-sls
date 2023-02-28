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
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"io"
	"net/http"
	"time"

	"github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	// 注册一个Metric指标（非必要步骤）
	labels := []attribute.KeyValue{
		attribute.String("label1", "value1"),
	}
	meter := global.Meter("aliyun.sls")
	sayDavidCount, _ := meter.Int64Counter("say_david_count")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if time.Now().Unix()%10 == 0 {
			_, _ = io.WriteString(w, "Hello, world!\n")
		} else {
			// 如果需要记录一些事件，可以获取Context中的span并添加Event（非必要步骤）
			ctx := req.Context()
			span := trace.SpanFromContext(ctx)
			span.AddEvent("say : Hello, I am david", trace.WithAttributes(attribute.KeyValue{
				Key:   "label-key-1",
				Value: attribute.StringValue("label-value-1"),
			}))

			_, _ = io.WriteString(w, "Hello, I am david!\n")
			sayDavidCount.Add(req.Context(), 1, labels...)
		}
	}

	// 使用 otel net/http的自动注入方式，只需要使用otelhttp.NewHandler包裹http.Handler即可
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	http.Handle("/hello", otelHandler)
	fmt.Println("Now listen port 8080, you can visit 127.0.0.1:8080/hello .")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
