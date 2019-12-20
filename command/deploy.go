package command

import (
	"errors"
	"fmt"
)

type Deploy struct {
	product string
}

func NewDeploy(args []string) (*Deploy, error) {
	if args == nil || len(args) == 0 {
		return nil, errors.New("Deploy requires 1 argument")
	}

	return &Deploy{product: args[0]}, nil
}

func (d Deploy) Execute(out chan string) {
	defer close(out)

	out <- fmt.Sprintf("Deploying *%s*...", d.product)
}
