package jaeger

import "github.com/opentracing/opentracing-go"

// JaegerAdapter is an adapter that wraps the Jaeger tracer to implement the JaegerTracer interface.
type JaegerAdapter struct {
	tracer opentracing.Tracer
}

// StartSpan starts a new span with the given operation name and options.
func (ja *JaegerAdapter) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return ja.tracer.StartSpan(operationName, opts...)
}

// Inject injects the given span context into the carrier.
func (ja *JaegerAdapter) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return ja.tracer.Inject(sm, format, carrier)
}

// Extract extracts a span context from the carrier.
func (ja *JaegerAdapter) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return ja.tracer.Extract(format, carrier)
}
