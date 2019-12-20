package slack

import (
	"kubot/command"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	server := newMockServer()

	startOptions = append(startOptions, slack.OptionAPIURL(server.httpURL()))
	rtm = slack.New("", startOptions...).NewRTM()

	go Start()

	rtm.IncomingEvents <- newMessageEvent("!help")
	rsp, err := server.waitForRequest()
	if err != nil {
		t.Fatal(err)
	}
	msg := getOutgoingMessage(rsp)

	assert.Equal(t, command.Help{}.Usage(), msg.Text)
}
