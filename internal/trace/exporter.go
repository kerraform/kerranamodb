package trace

import (
	"io"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

type ExporterType string

const (
	ExporterTypeConsole ExporterType = "console"
	ExporterTypeJaeger  ExporterType = "jaeger"
)

func NewConsoleExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

func NewJaegerExporter(url string) (trace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
		),
	)
}
