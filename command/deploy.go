package command

import (
	"fmt"

	"github.com/apex/log"
)

type Deploy struct {
	product string
	version string
	release string
}

func NewDeploy(args []string) (*Deploy, error) {
	if args == nil || len(args) < 2 {
		return nil, &CommandArgumentError{"usage: deploy <product> <version> [release]"}
	}

	deploy := &Deploy{
		product: args[0],
		version: args[1],
		release: args[0],
	}

	if len(args) > 2 {
		deploy.release = args[2]
	}

	return deploy, nil
}

func (d Deploy) Execute(out chan string, context Context) {
	defer close(out)

	out <- fmt.Sprintf("Deploying *%s-%s* to *%s* environment...", d.product, d.version, context.Environment.Name)

	if err := Execute("deploy", d.product, map[string]string{"product": d.product, "version": d.version, "release": d.release, "environment": context.Environment.Name}, context.Environment.Variables); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* deployment failed", d.product)
		return
	}
	out <- fmt.Sprintf("*%s-%s* was successfully deployed.", d.product, d.version)
}
