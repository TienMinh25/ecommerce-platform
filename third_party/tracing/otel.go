package tracing

import (
	"context"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

// NewTracer create new function in service main and inject it into its service
func NewTracer(env *env.EnvManager, serviceName string) (pkg.Tracer, error) {
	grpcExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(env.Jeager.JeagerExporterGrpcEndpoint)),
	)

	if err != nil {
		return nil, errors.Wrap(err, "otlptrace.New")
	}

	// create resource for analyze trace
	resourceTrace, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceVersion("v1"),
			semconv.ServiceName(serviceName),
			attribute.String("library.language", "go"),
		),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Faield to create a resource")
	}

	// create tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			grpcExporter,
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize)),
		sdktrace.WithResource(resourceTrace),
	)

	// register with otel tracer provider
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &tracer{
		provider: tracerProvider,
		tracer:   tracerProvider.Tracer(serviceName),
	}, nil
}

func (t *tracer) StartFromContext(ctx context.Context, name string) (context.Context, pkg.Span) {
	ctx, s := t.tracer.Start(ctx, name)

	return ctx, &span{
		span: s,
	}
}

func (t *tracer) StartFromSpan(ctx context.Context, span pkg.Span, name string) (context.Context, pkg.Span) {
	return t.StartFromContext(span.Context(ctx), name)
}

func (t *tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

func (t *tracer) Inject(ctx context.Context, carrier pkg.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func (t *tracer) Extract(ctx context.Context, carrier pkg.TextMapCarrier) pkg.Span {
	return &span{
		span: trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, carrier)),
	}
}

type span struct {
	span trace.Span
}

func (s *span) RecordError(err error) {
	s.span.RecordError(err)
}

func (s *span) SetAttributes(key string, value string) {
	s.span.SetAttributes(attribute.String(key, value))
}

func (s *span) Context(ctx context.Context) context.Context {
	return trace.ContextWithSpan(ctx, s.span)
}

func (s *span) End() {
	s.span.End()
}
