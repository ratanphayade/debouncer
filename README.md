# Debouncer

Debounce and throttle are two similar (but different!) techniques to control how many times we allow a function to be
executed over time. Debounce is used when you have to consider only the final state.

A typical example for this is auto submit / auto suggestions. In these cases rather than making a request to server for
suggestions, it's recommended to wait for some time. If the end user is still typing then ignore the previous request
and consider the latest value.

This can be very well adapted in multiple places to reduce the load on the system.

## Getting Started
install `debouncer`
```
$ go get -u -v github.com/ratanphayade/debouncer
```

### Using Debouncer

#### Create an instance
To create an instance
```go
d := debouncer.NewDebouncer(ttl)
```
Here ttl is the interval until which the event has to wait before performing action. In case before ttl another event
received then, previous event will be ignored, and it'll again wait for a duration specified by ttl before performing the
action

*Note: One instance can handle only action. If required to have multiple action then, multiple instance has to be created*

 #### Start the action listener
 Action should be of type
```go
type Action func(ctx context.Context, value interface{}) error
```

Once we have the action defined, we can start the action listener
```go
d.Do(ctx, func(_ context.Context, val interface{}) error {
    counter++
    result = result + val.(int)
    return nil
})
```

#### Triggering the event
Calling below method with every event will invoke the debounce action
```go
d.TriggerAction(value)
```

