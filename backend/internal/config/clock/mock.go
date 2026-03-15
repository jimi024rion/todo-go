package clock

import (
	"context"
	"time"
)

// MockClock is a mock implementation of Clock for testing.
type MockClock struct {
	FixedTime time.Time
}

// NewMockClock returns a new MockClock with the given fixed time.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{
		FixedTime: t,
	}
}

// Now returns the time defined in FixedTime, truncated to seconds.
func (c *MockClock) Now(ctx context.Context) time.Time {
	return c.FixedTime.Truncate(time.Second)
}
