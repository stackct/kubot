package command

import (
	"fmt"
	"github.com/apex/log"
)

type Deploy struct {
	product string
	version string
}

func NewDeploy(args []string) (*Deploy, error) {
	if args == nil || len(args) != 2 {
		return nil, &CommandArgumentError{"Deploy requires 2 arguments"}
	}
	return &Deploy{product: args[0], version: args[1]}, nil
}

func (d Deploy) Execute(out chan string, context Context) {
	defer close(out)

	out <- fmt.Sprintf("Deploying *%s-%s* to *%s* environment...", d.product, d.version, context.Environment.Name)

	if err := Execute("deploy", map[string]string{"product": d.product, "release": context.Environment.Release, "version": d.version, "environment": context.Environment.Name}); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* deployment failed", d.product)
		return
	}
	out <- fmt.Sprintf("*%s-%s* was successfully deployed.", d.product, d.version)
}
