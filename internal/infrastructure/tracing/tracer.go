package tracing

import (
	"github.com/opentracing/opentracing-go"
)

// Tracer is an interface that defines the tracing methods required by the application.
type Tracer interface {
	StartSpan(operationName string, opts ...interface{}) interface{}
	Inject(sm interface{}, format interface{}, carrier interface{}) error
	Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error)
}
