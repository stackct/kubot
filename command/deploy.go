package command

import (
	"errors"
	"fmt"
)

type Deploy struct {
	product string
	version string
}

func NewDeploy(args []string) (*Deploy, error) {
	if args == nil || len(args) != 2 {
		return nil, errors.New("Deploy requires 2 arguments")
	}
	return &Deploy{product: args[0], version: args[1]}, nil
}

func (d Deploy) Execute(out chan string) {
	defer close(out)

	out <- fmt.Sprintf("Deploying *%s*...", d.product)
	if err := Execute("deploy", map[string]string{"product": d.product, "version": d.version}); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* deployment failed", d.product)
		return
	}
	out <- fmt.Sprintf("*%s* deployment complete", d.product)
}
