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

func TestPassingContext(t *testing.T) {
	var state1, state2, state3 machine.State

	var result int

	state1 = func(runtime machine.Runtime) {
		first := 1
		ctx := context.WithValue(runtime.Context(), "first", first)
		runtime.NextState(ctx, state2)
	}

	state2 = func(runtime machine.Runtime) {
		first := runtime.Context().Value("first").(int)
		ctx := context.WithValue(runtime.Context(), "second", first+1)
		runtime.NextState(ctx, state3)
	}

	state3 = func(runtime machine.Runtime) {
		second := runtime.Context().Value("second").(int)
		result = second + 1
		runtime.NextState(runtime.Context(), nil)
	}

	runtime := local.Runtime(context.Background(), state1)

	<-runtime.Context().Done()

	if result != 3 {
		t.Errorf("expected %d, but got %d", 3, result)
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
