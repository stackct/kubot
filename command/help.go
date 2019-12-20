package command

type Help struct{}

func NewHelp() (*Help, error) {
	return &Help{}, nil
}

func (h Help) Name() string {
	return "help"
}

func (h Help) Execute(out chan string) {
	defer close(out)

	out <- h.Usage()
}

func (h Help) Usage() string {
	return "available commands: [deploy]"
}
