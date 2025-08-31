package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" // Add this
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func newTracerProvider() (*trace.TracerProvider, error) {
	// Console exporter for debugging
	consoleExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// OTLP exporter to send traces to collector
	otlpExporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cineplex"),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		// Console exporter - prints to stdout for debugging
		trace.WithBatcher(consoleExporter, trace.WithBatchTimeout(time.Second)),
		// OTLP exporter - sends to collector
		trace.WithBatcher(otlpExporter, trace.WithBatchTimeout(time.Second*5)),
	)

	return tracerProvider, nil
}

//func newTracerProvider() (*trace.TracerProvider, error) {
//	traceExporter, err := stdouttrace.New(
//		stdouttrace.WithPrettyPrint())
//	if err != nil {
//		return nil, err
//	}
//
//	tracerProvider := trace.NewTracerProvider(
//		trace.WithBatcher(traceExporter,
//			// Default is 5s. Set to 1s for demonstrative purposes.
//			trace.WithBatchTimeout(time.Second)),
//	)
//
//	return tracerProvider, nil
//}
//
