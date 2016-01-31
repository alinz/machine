package machine

import (
	"sync"

	"golang.org/x/net/context"
)

//State this is the state fucntion signiture.
type State func(Runtime)

//Runtime is base blocking start point.
type Runtime interface {
	//NextState by calling NextState, State Machine goes to another state.
	//one thing to remember, once you call this method, you are no longer
	//in that state. Make sure that you are nor doing anything right after this
	//call.
	//Passing nil to State tells the Runtime that this is the end of state machine.
	NextState(context.Context, State)
	//return the current context of Runtime.
	Context() context.Context
	//Fork run each state as initial state of unique Runtime.
	//I does not block, if you want to block until all of them are done,
	//pass the array of runtime to Join function.
	Fork(context.Context, ...State) []Runtime
}

//Join It blocks until all Runtimes finishes their job
func Join(runtimes []Runtime) {
	var wg sync.WaitGroup
	wg.Add(len(runtimes))

	for _, runtime := range runtimes {
		go func(runtime Runtime) {
			defer wg.Done()
			<-runtime.
				Context().
				Done()
		}(runtime)
	}

	wg.Wait()
}
