package deferrer

type Deferrer struct {
	deferred []func()
}

func (d *Deferrer) RunDeferred() {
	for _, fn := range d.deferred {
		fn()
	}
}

func (d *Deferrer) Defer(fn func()) {
	d.deferred = append(d.deferred, fn)
}
