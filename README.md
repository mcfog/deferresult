# deferresult

## Motivation

Unhandled error in defer statement is always a problem. We have some wrapper utilities in our project to eat errors with Logger looks like this:

```go
func (l *logger) Wrap(f func() error, msg string) func() {
	return func() {
		err := f()
		// log the err down
    }
}
```

But that caused another trap: `defer logger.Wrap(resource.Close, "closing resource")` would NOT close the resource since the `func()` returned by `Wrap` is not called!

Basically we found that expression in defer should never return anything so this linter is created to check about this.
