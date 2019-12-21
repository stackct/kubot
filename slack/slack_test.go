package slack

import (
	"errors"
	"kubot/command"
	"kubot/config"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

var server *MockServer

func init() {
	server = newMockServer()
}

func BeforeEach(t *testing.T) func() {
	Conf = config.NewMockConfig()

	users = []slack.User{
		slack.User{ID: "_user", Profile: slack.UserProfile{Email: "foo@invalid.co"}},
	}
	startOptions = append(startOptions, slack.OptionAPIURL(server.httpURL()))
	rtm = slack.New("", startOptions...).NewRTM()

	go Start()

	return Teardown
}

func Teardown() {
	parser = command.NewSlackCommandParser()
}

func TestStart_Invalid_Command(t *testing.T) {
	defer BeforeEach(t)()

	testCases := []struct {
		option       command.MockParserOption
		incomingText string
		outgoingText string
	}{
		{ParseError, "any", "unknown command"},
		{ParseSuccess, "any", "fin"},
	}

	for _, tc := range testCases {
		parser = command.NewMockParser(tc.option)
		runTest(t, tc.incomingText, tc.outgoingText)
	}
}

func TestStart_Wrong_Environment(t *testing.T) {
	defer BeforeEach(t)()

	testCases := []struct {
		option       config.MockConfigOption
		incomingText string
		outgoingText string
	}{
		{MissingEnvironment, "!help", "Environment not found"},
		{ExistingEnvironment, "!help", "available commands: [deploy]"},
	}

	for _, tc := range testCases {
		Conf = config.NewMockConfig(tc.option)

		runTest(t, tc.incomingText, tc.outgoingText)
	}
}

func TestStart_No_Access(t *testing.T) {
	defer BeforeEach(t)()

	testCases := []struct {
		option       config.MockConfigOption
		incomingText string
		outgoingText string
	}{
		{NoAccess, "!help", "Authorization denied"},
		{FullAccess, "!help", "available commands: [deploy]"},
	}

	for _, tc := range testCases {
		Conf = config.NewMockConfig(tc.option)

		runTest(t, tc.incomingText, tc.outgoingText)
	}
}

func ParseSuccess(mp *command.MockParser) {
	mp.MockError = nil
	mp.MockCommand = command.MockCommand{}
}

func ParseError(mp *command.MockParser) {
	mp.MockError = errors.New("unknown command")
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

func runTest(t *testing.T, incomingText string, outgoingText string) {
	rtm.IncomingEvents <- newMessageEvent(incomingText)
	rsp, err := server.waitForRequest()
	if err != nil {
		t.Fatal(err)
	}
	msg := getOutgoingMessage(rsp)

	assert.Equal(t, outgoingText, msg.Text)
}
