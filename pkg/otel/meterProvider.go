package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func newMeterProvider() (*metric.MeterProvider, error) {
	// OTLP exporter - will automatically use OTEL_EXPORTER_OTLP_ENDPOINT
	otlpExporter, err := otlpmetrichttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	// Optional: Console exporter for debugging
	consoleExporter, err := stdoutmetric.New()
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

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		// OTLP exporter - sends to collector using environment variables
		metric.WithReader(metric.NewPeriodicReader(
			otlpExporter,
			metric.WithInterval(time.Second*10),
		)),
		// Console exporter - for debugging (optional)
		metric.WithReader(metric.NewPeriodicReader(
			consoleExporter,
			metric.WithInterval(time.Second*5),
		)),
	)

	return meterProvider, nil
}
