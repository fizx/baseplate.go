package metricsbp_test

import (
	"github.com/fizx/baseplate.go/metricsbp"
	"github.com/fizx/baseplate.go/tracing"
)

// This example demonstrates how to use CreateServerSpanHook.
func ExampleCreateServerSpanHook() {
	const prefix = "service.server"

	// initialize the CreateServerSpanHook
	hook := metricsbp.CreateServerSpanHook{}

	// register the hook with Baseplate
	tracing.RegisterCreateServerSpanHooks(hook)
}
