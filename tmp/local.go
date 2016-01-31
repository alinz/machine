package Machine

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

type localDone chan struct{}

func (l localDone) Wait(timeout int64) {
	if timeout > 0 {
		select {
		case <-l:
		case <-time.After(time.Duration(timeout)):
		}
	} else {
		<-l
	}
}

type localState struct {
	Context context.Context
	State   State
}
type localStateTransition chan localState

func (ls localStateTransition) Next(ctx context.Context, state State) {
	ls <- localState{ctx, state}
}

func (ls localStateTransition) Fork(ctx context.Context, states ...State) Done {
	var wg sync.WaitGroup
	wg.Add(len(states))

	var done localDone
	done = make(chan struct{})

	for _, state := range states {
		go func(initialState State) {
			defer wg.Done()

			NewLocalMachine().
				RunStateMachine(ctx, initialState).
				Wait(0)
		}(state)
	}

	go func() {
		defer close(done)
		wg.Wait()
	}()

	return done
}

func (ls localStateTransition) Done() {
	close(ls)
}

type localMachine struct {
	done         localDone
	transitioner localStateTransition
}

func (lm *localMachine) RunStateMachine(ctx context.Context, initialState State) Done {
	go func() {
		var localState localState
		ok := true

		for ok {
			select {
			case localState, ok = <-lm.transitioner:
				if ok {
					localState.State(localState.Context, lm.transitioner)
				}
			case _, ok = <-ctx.Done():
			}
		}

		defer close(lm.done)
	}()

	initialState(ctx, lm.transitioner)

	return lm.done
}

//NewLocalMachine is a simple implemenation of state Machine which uses go
//channel
func NewLocalMachine() Machine {
	return &localMachine{
		done:         make(chan struct{}),
		transitioner: make(chan localState, 1),
	}
}
