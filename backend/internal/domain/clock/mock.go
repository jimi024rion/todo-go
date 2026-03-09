package clock

import (
	"context"
	"time"
)

// MockClock is a mock implementation of Clock for testing.
type MockClock struct {
	NowFunc func(ctx context.Context) time.Time
}

// NewMockClock returns a new MockClock with the given fixed time.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{
		NowFunc: func(ctx context.Context) time.Time {
			return t
		},
	}
}

// Now returns the time defined in NowFunc.
func (c *MockClock) Now(ctx context.Context) time.Time {
	if c.NowFunc != nil {
		return c.NowFunc(ctx)
	}
	return time.Time{}
}
