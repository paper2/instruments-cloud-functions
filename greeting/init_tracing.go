package greeting

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracing() *sdk.TracerProvider {
	var opts []otlptracehttp.Option
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	client := otlptracehttp.NewClient(opts...)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithAttributes(
		// Cusom attributes
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	// W3C Trace Context propagator
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
