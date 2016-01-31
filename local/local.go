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

func (l *localRuntime) loop(start machine.State) {
	go func() {
		var localItem *item
		ok := true

		for ok {
			select {
			case localItem, ok = <-l.local:
				if ok {
					localItem.state(l)
				}
			case _, ok = <-l.ctx.Done():
			}
		}

		defer close(l.local)
	}()

	start(l)
}

func LocalRuntime(ctx context.Context, initialState machine.State) machine.Runtime {
	runtime := localRuntime{
		local: make(chan *item, 1),
		ctx:   ctx,
	}

	runtime.loop(initialState)

	return &runtime
}
