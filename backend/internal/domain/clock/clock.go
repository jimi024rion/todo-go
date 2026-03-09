package clock

import (
	"context"
	"time"
)

// Clock is an interface for time operations.
type Clock interface {
	Now(ctx context.Context) time.Time
}
