package txnlog

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey struct{}

// loggerKey is the context key for storing the Logger.
var loggerKey = contextKey{}

// WithLogger returns a new context with the Logger stored in it.
func WithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext extracts the Logger from the context.
// Returns nil if no Logger is present.
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerKey).(*Logger); ok {
		return l
	}
	return nil
}

// MustFromContext extracts the Logger from the context.
// Panics if no Logger is present. Use only when a Logger is guaranteed.
func MustFromContext(ctx context.Context) *Logger {
	l := FromContext(ctx)
	if l == nil {
		panic("txnlog: no Logger in context")
	}
	return l
}
