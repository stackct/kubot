package slack

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"kubot/command"
	"kubot/config"
)

var channel string

func init() {
	channel = "_channel"
}

func newMessageEvent(body string) slack.RTMEvent {
	return slack.RTMEvent{
		Type: "Message",
		Data: &slack.MessageEvent{
			Msg: slack.Msg{
				User:    "_user",
				Channel: channel,
				Text:    body,
			},
		},
	}
}

func getOutgoingMessage(s string) slack.OutgoingMessage {
	var resp slack.OutgoingMessage
	err := json.Unmarshal([]byte(s), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func ParseSuccess(mp *command.MockParser) {
	mp.MockError = nil
	mp.MockCommand = command.MockCommand{}
}

func ParseUnknownCommandError(mp *command.MockParser) {
	mp.MockError = command.NewUnknownCommandError("unknown command")
	mp.MockCommand = nil
}

func ParseCommandArgumentError(mp *command.MockParser) {
	mp.MockError = command.NewCommandArgumentError("command argument error")
	mp.MockCommand = nil
}

func MissingEnvironment(c *config.MockConfig) {
	c.MockHasEnvironment = false
}

func ExistingEnvironment(c *config.MockConfig) {
	c.MockHasEnvironment = true
}

func NoAccess(c *config.MockConfig) {
	c.MockHasAccess = false
}

func FullAccess(c *config.MockConfig) {
	c.MockHasAccess = true
}
