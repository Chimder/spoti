package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	// semconv "go.opentelemetry.io/otel/semconv@latest"
)

func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("spoti-api"),
			semconv.ServiceVersionKey.String("1.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, err
	}

	promExp, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExp),
	)
	otel.SetMeterProvider(meterProvider)

	// traceExp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	return nil, err
	// }
	// traceProvider := sdktrace.NewTracerProvider(
	// 	sdktrace.WithResource(res),
	// 	sdktrace.WithBatcher(traceExp),
	// 	sdktrace.WithSampler(sdktrace.AlwaysSample()),
	// )
	// otel.SetTracerProvider(traceProvider)

	shutdown = func(ctx context.Context) error {
		var errs []error
		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		// if err := traceProvider.Shutdown(ctx); err != nil {
		// 	errs = append(errs, err)
		// }
		if len(errs) > 0 {
			log.Printf("OTel shutdown: %v", errs)
		}
		return nil
	}

	log.Println("init OTel")
	return shutdown, nil
}
