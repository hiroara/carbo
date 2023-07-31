package task

import (
	"time"

	"github.com/hiroara/carbo/task/internal/inout"
)

// An option for a Task's input or output.
type InOutOption func(*inout.Options)

func newOptions(options []InOutOption) *inout.Options {
	opts := &inout.Options{}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

// Return an option to configure timeout.
// This timeout is applied to sending or receiving a value to/from a channel.
// Please note that this will not be applied to sending or receiving an element rather than an entire Task execution.
// For a Task input, when receiving a value from the input channel takes more than the timeout value, the input channel
// will be closed, and the passed context will be canceled.
// For a Task output, when sending a value to the output channel takes more than the timeout value, the context
// passed to the task will be canceled.
func WithTimeout(d time.Duration) InOutOption {
	return func(opts *inout.Options) {
		opts.Timeout = d
	}
}

type options struct {
	inOpts  []InOutOption
	outOpts []InOutOption
}

// An option for a task.
type Option func(opts *options)

// An Option to set input options.
func WithInputOptions(opts ...InOutOption) Option {
	return func(tOpts *options) {
		tOpts.inOpts = opts
	}
}

// An Option to set output options.
func WithOutputOptions(opts ...InOutOption) Option {
	return func(tOpts *options) {
		tOpts.outOpts = opts
	}
}
