package local

import (
	"github.com/alinz/machine"
	"golang.org/x/net/context"
)

type item struct {
	state machine.State
	ctx   context.Context
}

type localRuntime struct {
	local chan *item
	ctx   context.Context
}

func (l *localRuntime) NextState(ctx context.Context, state machine.State) {
	l.local <- &item{state, ctx}
}

func (l *localRuntime) Context() context.Context {
	return l.ctx
}

func (l *localRuntime) Fork(ctx context.Context, states ...machine.State) []machine.Runtime {
	var runtimes []machine.Runtime
	runtimes = make([]machine.Runtime, 0)

	for _, state := range states {
		runtimes = append(runtimes, Runtime(ctx, state))
	}

	return runtimes
}

func (l *localRuntime) loop(start machine.State) {
	var cancel context.CancelFunc
	l.ctx, cancel = context.WithCancel(l.ctx)

	go func() {
		var localItem *item
		ok := true

		for ok {
			select {
			case localItem, ok = <-l.local:
				if ok {
					if localItem.state != nil {
						localItem.state(l)
					} else {
						//we are going to cancel if state sends a nil state as a next state
						//this implies that state reaches the end.
						cancel()
					}
				}
			case _, ok = <-l.ctx.Done():
			}
		}

		defer close(l.local)
	}()

	start(l)
}

//Runtime creates a Runtime based on go channel.
func Runtime(ctx context.Context, initialState machine.State) machine.Runtime {
	runtime := localRuntime{
		local: make(chan *item, 1),
		ctx:   ctx,
	}

	runtime.loop(initialState)

	return &runtime
}
