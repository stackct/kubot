package command

import (
	"errors"
	"fmt"
	"io"
	"kubot/config"
	"os/exec"
	"regexp"
	"strings"
)

type Command interface {
	Name() string
	Execute(output chan string)
}

type SlackCommandParser struct{}

var log = config.Log

func NewSlackCommandParser() SlackCommandParser {
	return SlackCommandParser{}
}

func (foo SlackCommandParser) Parse(c string) (Command, error) {
	re, _ := regexp.Compile(`^!(?P<command>help|deploy) ?(?P<args>.*)?`)

	keys := re.SubexpNames()
	vals := re.FindAllStringSubmatch(c, -1)

	if len(vals) == 0 {
		return nil, errors.New("unknown command")
	}

	md := map[string]string{}
	for i, n := range vals[0] {
		md[keys[i]] = n
	}

	args := strings.Fields(md["args"])

	switch md["command"] {
	case "help":
		return NewHelp()
	case "deploy":
		return NewDeploy(args)
	}

	return nil, nil
}

func Execute(c Command, writer *io.PipeWriter, out chan string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = writer
	cmd.Stderr = writer

	if err := cmd.Start(); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* command failed", c.Name())
		return
	}
	if err := cmd.Wait(); err != nil {
		log.Error(err.Error())
		out <- fmt.Sprintf("*%s* command failed", c.Name())
		return
	}
}
