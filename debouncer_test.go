package debouncer_test

import (
	"context"
	"testing"
	"time"

	"github.com/ratanphayade/debouncer"
	"github.com/stretchr/testify/assert"
)

func TestDebounce_Do(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize the debouncer with the interval
	// until which the trigger will wait before executing the action
	d := debouncer.NewDebouncer(100 * time.Millisecond)
	var (
		counter int
		result  int
	)

	// start the action listener
	d.Do(ctx, func(_ context.Context, val interface{}) error {
		counter++
		result = result + val.(int)
		return nil
	})

	// triggering multiple action events
	for i := 0; i < 5; i++ {
		d.TriggerAction(i)
	}

	time.Sleep(200 * time.Millisecond)
	cancel()

	assert.Equal(t, 1, counter)
	assert.Equal(t, 4, result)
}
