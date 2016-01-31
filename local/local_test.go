package local_test

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/alinz/machine"
	"github.com/alinz/machine/local"
)

func TestLocalRuntime(t *testing.T) {

	var state1, state2, state3 machine.State

	state1 = func(runtime machine.Runtime) {
		fmt.Println("Hello from stage1")
		runtime.NextState(runtime.Context(), state2)
	}

	state2 = func(runtime machine.Runtime) {
		fmt.Println("Hello from stage2")
		runtime.NextState(runtime.Context(), state3)
	}

	state3 = func(runtime machine.Runtime) {
		fmt.Println("Hello from stage3")
		runtime.NextState(runtime.Context(), nil)
	}

	runtime := local.LocalRuntime(context.Background(), state1)

	<-runtime.Context().Done()
}
