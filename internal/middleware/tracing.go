package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TraceRequests middleware traces each incoming HTTP request
func TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("http-tracer")
		ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
		span.SetAttributes(attribute.String("http.method", r.Method))
		span.SetAttributes(attribute.String("http.url", r.URL.String()))

		defer span.End()

		// Pass the context with the span to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
