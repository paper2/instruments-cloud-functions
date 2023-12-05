package greeting

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Flush interface {
	ForceFlush(context.Context) error
}

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

// InstrumentedHandler wraps the function with OpenTelemetry instrumentation.
func InstrumentedHandler(functionName string, function HttpHandler, flusher Flush) HttpHandler {
	opts := []trace.SpanStartOption{
		// customizable span attributes
		trace.WithAttributes(semconv.FaaSTriggerHTTP),
	}

	// create instrumented handler
	handler := otelhttp.NewHandler(
		http.HandlerFunc(function), functionName, otelhttp.WithSpanOptions(opts...),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		// call the actual handler
		handler.ServeHTTP(w, r)

		// NOTE: ForceFlush() may extend the function's duration. It must be used carefully.
		//       If ForceFlush() is not called, spans are send on background.
		//       Backgraound tasks are not recommended in Cloud Functions. Span data sometimes get lost.
		err := flusher.ForceFlush(r.Context())
		if err != nil {
			log.Printf("failed to flush spans: %v", err)
		}
	}
}
