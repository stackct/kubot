package command

import (
	"fmt"
	"kubot/config"
)

type Help struct{}

func NewHelp() (*Help, error) {
	return &Help{}, nil
}

func (h Help) Execute(out chan string, context Context) {
	defer close(out)

	out <- h.Usage()
}

func (h Help) Usage() string {
	return fmt.Sprintf("available commands: %v", config.Conf.GetCommands())
}
