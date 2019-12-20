package slack

import (
	"encoding/json"
	"github.com/nlopes/slack"
)

func newMessageEvent(body string) slack.RTMEvent {
	return slack.RTMEvent{
		Type: "Message",
		Data: &slack.MessageEvent{
			Msg: slack.Msg{
				User:    "_user",
				Channel: "_channel",
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
