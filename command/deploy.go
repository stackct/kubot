package command

import (
	"errors"
	"fmt"
)

type Deploy struct {
	product   string
	version   string
	repo      string
	timeout   string
	chartFile string
}

func NewDeploy(args []string) (*Deploy, error) {
	if args == nil || len(args) == 0 {
		return nil, errors.New("Deploy requires 1 argument")
	}
	return &Deploy{product: args[0]}, nil
}

func (d Deploy) Name() string {
	return "deploy"
}

func (d Deploy) Execute(out chan string) {
	defer close(out)

	//writer := log.StandardLogger().Writer()
	out <- fmt.Sprintf("Deploying *%s*...", d.product)
	//Execute(d, out, "helm", "repo", "update")
	//Execute(d, out, "helm", "upgrade", d.product, d.repo+"/"+d.product, "--wait", "--timeout", d.timeout, "--version", d.version, "-f", d.chartFile, "--set", d.product+".image.tag="+d.version)
	out <- fmt.Sprintf("*%s* deployment complete", d.product)
}
