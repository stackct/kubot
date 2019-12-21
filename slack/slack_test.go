package slack

import (
	"kubot/command"
	"kubot/config"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func init() {
	users = []slack.User{
		slack.User{ID: "_user", Profile: slack.UserProfile{Email: "foo@invalid.co"}},
	}
}

func TestStart(t *testing.T) {
	Conf = MockConfig{}

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

type MockConfig struct{}

func (c MockConfig) HasAccess(s string, e string) bool {
	return true
}

func (c MockConfig) GetEnvironmentByChannel(ch string) (*config.Environment, error) {
	return &config.Environment{Name: ch, Channel: ch}, nil
}
