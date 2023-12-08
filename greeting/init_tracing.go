package greeting

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracing initializes tracing.
func InitTracing() *sdk.TracerProvider {
	var opts []otlptracehttp.Option
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		// In local environment, TLS is not set up.
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithAttributes(
		// Custom attributes
		),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)

	// Set the global TraceProvider to the SDK's TraceProvider.
	otel.SetTracerProvider(tp)

	// W3C Trace Context propagator
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
