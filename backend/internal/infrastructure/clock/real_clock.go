package clock

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/clock"
)

type realClock struct{}

// NewRealClock returns a new RealClock implementation.
func NewRealClock() clock.Clock {
	return &realClock{}
}

// Now returns the current time in the location specified in the context.
func (c *realClock) Now(ctx context.Context) time.Time {
	loc := clock.LocationFromContext(ctx)
	return time.Now().In(loc)
}
