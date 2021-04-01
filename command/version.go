package command

import (
	"github.com/apex/log"
)

type Version struct {
	product string
}

func NewVersion(args []string) (*Version, error) {
	if args == nil || len(args) < 1 {
		return nil, &CommandArgumentError{"usage: version <product>"}
	}

	version := &Version{
		product: args[0],
	}

	return version, nil
}

func (v Version) Execute(out chan string, context Context) {
	defer close(out)

	if err := Execute("version", v.product, map[string]string{"version": v.product, "environment": context.Environment.Name}, context.Environment.Variables, out); err != nil {
		log.Error(err.Error())
		out <- "failed to execute version command"
		return
	}
}
