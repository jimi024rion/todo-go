package clock

import (
	"context"
	"time"
)

type contextKey string

const (
	locationKey contextKey = "timezone_location"
)

// ContextWithLocation returns a new context with the given location.
func ContextWithLocation(ctx context.Context, loc *time.Location) context.Context {
	return context.WithValue(ctx, locationKey, loc)
}

// LocationFromContext returns the location from the context.
// If not set, it returns time.UTC as default.
func LocationFromContext(ctx context.Context) *time.Location {
	if loc, ok := ctx.Value(locationKey).(*time.Location); ok {
		return loc
	}
	return time.UTC
}
