# tracer

A minimalist mechanism to support adding a `traceid` and optional
`traceorigin` to context, as well as a `logr.Logger` also containing
these values.

The intended use case is where you need a context and logger with a
single value which is used for tracing purposes across multiple log
entries.

## API stability

The API is not stable.  To avoid problems use tagged releases.

Whilst this code is used by several projects, I've changed the API
more than once and am not yet happy with it.
