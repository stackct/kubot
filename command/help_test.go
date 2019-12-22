package command

import (
	"kubot/config"
	"testing"
)

func TestHelpExecute(t *testing.T) {
	out := make(chan string)
	config.Conf = config.NewMockConfig()

	expected := "available commands: [deploy]"

	go func() {
		if msg := <-out; expected != msg {
			t.Fatalf("Assertion failed; wanted: %s, but got: %s", expected, msg)
		}
	}()

	Help{}.Execute(out)
}
