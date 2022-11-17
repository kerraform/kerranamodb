package trace

import (
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewTracer(res *resource.Resource, exp trace.SpanExporter) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
}
