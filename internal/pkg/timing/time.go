package timing

import "time"

// Interface Timing describes custom time provider.
// It can be very useful for testing
type Timing interface {
	Now() time.Time
	Since(start time.Time) time.Duration
}

// Struct timingProvider implements Timing
type timingProvider struct{}

func NewTimingProvider() *timingProvider {
	return &timingProvider{}
}

// Now returns current time
func (t *timingProvider) Now() time.Time {
	return time.Now()
}

// Since returns duration since start
func (t *timingProvider) Since(start time.Time) time.Duration {
	return time.Since(start)
}
