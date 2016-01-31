package local_test

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/alinz/machine"
	"github.com/alinz/machine/local"
)

func TestRuntime(t *testing.T) {
	var state1, state2, state3 machine.State

	const finalAnswer = 10
	counter := 0

	state1 = func(runtime machine.Runtime) {
		counter += 5
		runtime.NextState(runtime.Context(), state2)
	}

	state2 = func(runtime machine.Runtime) {
		counter += 6
		runtime.NextState(runtime.Context(), state3)
	}

	state3 = func(runtime machine.Runtime) {
		counter--
		runtime.NextState(runtime.Context(), nil)
	}

	runtime := local.Runtime(context.Background(), state1)

	<-runtime.
		Context().
		Done()

	if counter != finalAnswer {
		t.Errorf("expected %d, but got %d", finalAnswer, counter)
	}
}

func TestRuntimeComplex(t *testing.T) {
	var state1, state2 machine.State

	state1 = func(runtime machine.Runtime) {
		var internalState1, internalState2 machine.State

		internalState1 = func(runtime machine.Runtime) {
			fmt.Println("internal state 1")
			runtime.NextState(runtime.Context(), internalState2)
		}

		internalState2 = func(runtime machine.Runtime) {
			fmt.Println("internal state 2")
			runtime.NextState(runtime.Context(), nil)
		}

		<-local.
			Runtime(context.Background(), internalState1).
			Context().
			Done()

		runtime.NextState(runtime.Context(), state2)
	}

	state2 = func(runtime machine.Runtime) {
		fmt.Println("state2")
		runtime.NextState(runtime.Context(), nil)
	}

	runtime := local.Runtime(context.Background(), state1)

	<-runtime.
		Context().
		Done()
}

//
//
//
