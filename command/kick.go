package command

import (
	"fmt"

	"github.com/apex/log"
)

type Kick struct {
	product string
}

func NewKick(args []string) (*Kick, error) {
	if args == nil || len(args) < 1 {
		return nil, &CommandArgumentError{"usage: kick <product>"}
	}

	kick := &Kick{
		product: args[0],
	}

	return kick, nil
}

func (k Kick) Execute(out chan string, context Context) {
	defer close(out)

	out <- fmt.Sprintf("Kicking *%s* in *%s* environment...", k.product, context.Environment.Name)

	if err := Execute("kick", k.product, map[string]string{"product": k.product, "environment": context.Environment.Name}, context.Environment.Variables, out); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* kick failed", k.product)
		return
	}
	out <- fmt.Sprintf("*%s* was successfully kicked.", k.product)
}
