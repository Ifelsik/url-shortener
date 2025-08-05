package timing

import "time"

// Interface Timing describes custom time provider.
// It can be very useful for testing
type Timing interface {
	Now() time.Time
}

// Struct timingProvider implements Timing
type timingProvider struct{}

func NewTimingProvider() *timingProvider {
	return &timingProvider{}
}

func (t *timingProvider) Now() time.Time {
	return time.Now()
}
