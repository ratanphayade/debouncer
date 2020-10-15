package debouncer

import (
	"context"
	"sync"
	"time"
)

// Action method signature of the action
// which has to be performed on event
type Action func(ctx context.Context, value interface{}) error

type Debouncer struct {
	// Input represents the change event
	Input chan interface{}

	// Interval represent the max time it should wait
	// before performing the Action
	Interval time.Duration

	// once will be used to ensure that the Do method
	// is called only once per deboubncer instance
	// as a single debounce can take care of only one operation
	// calling it multiple times might cause inconsistencies
	once sync.Once
}

// NewDebouncer creates a new instance of debouncer
// this will create an unbuffered channel to capture a event
func NewDebouncer(interval time.Duration) *Debouncer {
	return &Debouncer{
		Input:    make(chan interface{}),
		Interval: interval,
	}
}

// TriggerAction records an event to perform the Action provide
// this will add given value to the input channel as notification
// for debouncer
func (d *Debouncer) TriggerAction(val interface{}) {
	d.Input <- val
}

// Do will run the debounce in a go routine
// and it'll make sure that its been invoked only once
// as multiple action can not fall under same config
func (d *Debouncer) Do(ctx context.Context, action Action) {
	// ensure debouncing is started only once per instance
	d.once.Do(func() {
		go d.debounce(ctx, action)
	})
}

// debounce will make sure that the given action is not performed repeatedly
// if its triggered multiple times within a given interval
func (d *Debouncer) debounce(ctx context.Context, action Action) {
	var (
		// hasEvent represents of there is a valid event received
		// this will help avoid unnecessary triggering if the method
		hasEvent bool

		// value holds the latest input received
		value interface{}

		timer = time.NewTimer(d.Interval)
	)

	for {
		select {
		// if there is an event then reset the timer
		// and update the hasEvent to true representing
		// to trigger the function once the timer ends
		case value = <-d.Input:
			hasEvent = true
			timer.Reset(d.Interval)

		// if the timer ends there is a valid event
		// then call the Action provider
		case <-timer.C:
			if hasEvent {
				_ = action(ctx, value)
				hasEvent = false
			}

		// if the application is being terminated
		// then stop the debouncing
		case <-ctx.Done():
			return
		}
	}
}
