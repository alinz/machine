package machine

import "golang.org/x/net/context"

//State this is the state fucntion signiture.
type State func(Runtime)

//Runtime is base blocking start point.
type Runtime interface {
	NextState(context.Context, State)
	Context() context.Context
	Fork(context.Context, ...State) []Runtime
}

func Join(runtimes []Runtime) {

}
