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

func RemoveDirectMessages(c string, contextName string) (string, error) {
	var directMessageMentions []string
	var commandParts []string

	for _, v := range strings.Fields(c) {
		if strings.HasPrefix(v, "@") {
			directMessageMentions = append(directMessageMentions, strings.TrimPrefix(v, "@"))
		} else {
			commandParts = append(commandParts, v)
		}
	}

	if len(directMessageMentions) == 0 || contains(directMessageMentions, contextName) {
		return strings.Join(commandParts, " "), nil
	}

	return "", errors.New("command is not intended for this instance")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (_ SlackCommandParser) Parse(c string) (Command, error) {
	re, err := regexp.Compile(`^\` + SlackCommandPrefix + `(?P<command>[a-z]+) ?(?P<args>.*)?`)
	if err != nil {
		return nil, err
	}

	keys := re.SubexpNames()
	vals := re.FindAllStringSubmatch(c, -1)

	if len(vals) == 0 {
		return nil, &UnknownCommandError{"input does not match command syntax"}
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
	case "bounce":
		return NewBounce(args)
	case "kick":
		return NewKick(args)
	case "version":
		return NewVersion(args)
	}

	return nil, &UnknownCommandError{"unknown command"}
}
