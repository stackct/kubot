package command

import "github.com/apex/log"

type Bounce struct {
}

func NewBounce(args []string) (*Bounce, error) {
	return &Bounce{}, nil
}

func (b Bounce) Execute(out chan string, context Context) {
	defer close(out)

	if err := Execute("bounce", "", map[string]string{"environment": context.Environment.Name}, context.Environment.Variables); err != nil {
		log.Error(err.Error())
		out <- "Bounce failed"
		return
	} else {
		out <- "Bouncing now"
	}
}
