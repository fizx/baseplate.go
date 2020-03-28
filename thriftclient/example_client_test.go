package thriftclient_test

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/fizx/baseplate.go/internal/gen-go/reddit/baseplate"
	"github.com/fizx/baseplate.go/log"
	"github.com/fizx/baseplate.go/thriftclient"
	"github.com/fizx/baseplate.go/tracing"
)

func ExampleMonitoredClient() {
	// variables should be properly initialized in production code
	var (
		transport thrift.TTransport
		factory   thrift.TProtocolFactory
	)
	// Create an actual service client
	client := baseplate.NewBaseplateServiceClient(
		// Use MonitoredClient to wrap a standard thrift client
		thriftclient.NewMonitoredClientFromFactory(transport, factory),
	)
	// Create a context with a server span
	_, ctx := opentracing.StartSpanFromContext(
		context.Background(),
		"test",
		tracing.SpanTypeOption{Type: tracing.SpanTypeServer},
	)
	// Calls should be automatically wrapped using client spans
	healthy, err := client.IsHealthy(ctx)
	log.Debug("%v, %s", healthy, err)
}
