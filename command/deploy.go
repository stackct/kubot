package command

import (
	"fmt"
	"github.com/apex/log"
	"kubot/config"
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

	commandName := "deploy"
	_, err := config.Conf.GetCommand(commandName, d.product)
	if err != nil {
		switch err.(type) {
		case *config.ProhibitedCmdError:
			log.Info(fmt.Sprintf("Skipped %s command for product %s because it was found on the prohibited command list", commandName, d.product))
		default:
			log.Error(err.Error())
			out <- fmt.Sprintf("%s command for product %s was not found", commandName, d.product)
		}
		return
	}

	productLabel := fmt.Sprintf("*%s-%s*", d.product, d.version)
	if d.release != d.product {
		productLabel = fmt.Sprintf("*%s* with %s", d.release, productLabel)
	}

	out <- fmt.Sprintf("Deploying %s to *%s*...", productLabel, context.Environment.Name)

	if err := Execute(commandName, d.product, map[string]string{"product": d.product, "version": d.version, "release": d.release, "environment": context.Environment.Name}, context.Environment.Variables, out); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("%s deployment *FAILED* on *%s*", productLabel, context.Environment.Name)
		return
	}
	log.WithField("product", d.product).WithField("version", d.version).WithField("environment", context.Environment.Name).WithField("username", context.User).Info("deployed successfully")
	out <- fmt.Sprintf("%s was successfully deployed to *%s*", productLabel, context.Environment.Name)
}
