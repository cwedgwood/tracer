// a minimalist tracing logger attached to context.Context

package tracer

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

// When storing keys in a context, use an unexported type to avoid collisions and ensure accessor functions are used.
type traceString string

func (t traceString) String() string { return string(t) }

const (
	keytraceid      traceString = "traceid"
	keytraceorigin  traceString = "traceorigin"
	GenerateTraceId string      = ""
)

func newtraceid() string { return "tr-" + uuid.New().String() }

// ContextLoggerWithTraceId creates a new context with traceid and logr.Logger.
//
// If a non-empty traceid is provided that will be used, otherwise a random value is generated.
//
// If a non-empty traceorigin is specified it will also be present in the context and the logger.  This should be considered
// optional.
func ContextLoggerWithTraceId(parentContext context.Context, parentLogger logr.Logger, traceid, traceorigin string) context.Context {
	if traceid == "" || traceid == GenerateTraceId {
		traceid = newtraceid()
	}
	tracingContext := context.WithValue(parentContext, keytraceid, traceid)
	keysAndValues := []any{
		keytraceid.String(), traceid,
	}
	if traceorigin != "" {
		keysAndValues = append(keysAndValues, keytraceorigin.String(), traceorigin)
		tracingContext = context.WithValue(tracingContext, keytraceorigin, traceorigin)
	}
	tracingLogger := parentLogger.WithValues(keysAndValues...)
	return logr.NewContext(tracingContext, tracingLogger)
}

// Return traceid and traceorigin from the context, or empty strings when not present.
func TraceIdAndOrigin(tracingContext context.Context) (string, string) {
	var traceid, traceorigin string
	ta := tracingContext.Value(keytraceid)
	if v, ok := ta.(string); ok {
		traceid = v
	}
	ta = tracingContext.Value(keytraceorigin)
	if v, ok := ta.(string); ok {
		traceorigin = v
	}
	return traceid, traceorigin
}
