package deferrer

// An embeddable struct that can be used to implement Defer on a struct.
type Deferrer struct {
	deferred []func()
}

// Run all registered deferred functions.
func (d *Deferrer) RunDeferred() {
	for _, fn := range d.deferred {
		fn()
	}
}

// Register a function to defer the execution.
func (d *Deferrer) Defer(fn func()) {
	d.deferred = append(d.deferred, fn)
}
