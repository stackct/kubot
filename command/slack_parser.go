package command

import (
	"errors"
	"kubot/config"
	"regexp"
	"strings"
)

var SlackCommandPrefix string

type SlackCommandParser struct{}

func init() {
	SlackCommandPrefix = config.Conf.GetCommandPrefix()
}

func NewSlackCommandParser() SlackCommandParser {
	return SlackCommandParser{}
}

func (foo SlackCommandParser) Parse(c string) (Command, error) {
	re, err := regexp.Compile(`^\` + SlackCommandPrefix + `(?P<command>[a-z]+) ?(?P<args>.*)?`)
	if err != nil {
		return nil, err
	}

	keys := re.SubexpNames()
	vals := re.FindAllStringSubmatch(c, -1)

	if len(vals) == 0 {
		return nil, errors.New("input does not match command syntax")
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

	return nil, errors.New("unknown command")
}
