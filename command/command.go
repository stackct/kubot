package command

import (
	"errors"
	"regexp"
	"strings"
)

type Command interface {
	Execute(output chan string)
}

type CommandParser interface {
	Parse(string) (Command, error)
}

type SlackCommandParser struct{}

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
