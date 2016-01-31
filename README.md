# Machine

## introduction

Machine is a go library to write state Machine. The idea came to me when I watched [Rob Pike's Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE). So the basic idea is that you create a State Machine and start with initial state. Every state either generates a new state or stops the system. There are no errors coming from state Runtime. Since an error can lead to another state.

I tried to provide a simple and powerful interfaces so one can easily extend the state Runtime into next level.


## details

there are only 1 interface which help me to abstract the complexity of state Runtime.

```go
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
```

And also the main type of my state Runtime `State` which is a simple function that accepts `Runtime` as the only argument.

```go
type State func(Runtime)
```

by just having these 2 things, we can easily build any state Runtimes.

## Next

- implement a distributed version of state Runtime using this interface
- implement more examples
