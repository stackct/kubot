package command

import (
	"testing"
)

func TestHelpExecute(t *testing.T) {
	out := make(chan string)

	expected := "available commands: [deploy]"

	go func() {
		if msg := <-out; expected != msg {
			t.Fatalf("Assertion failed; wanted: %s, but got: %s", expected, msg)
		}
	}()

	Help{}.Execute(out)
}
