package slack

import (
	"kubot/command"
	"kubot/config"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

var server *MockServer

func init() {
	server = newMockServer()
	command.SlackCommandPrefix = "!"
}

func BeforeEach(t *testing.T) func() {
	config.Conf = config.NewMockConfig()

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
		{ParseError, "any", ""},
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
		{MissingEnvironment, "!cmd", "Environment not found"},
		{ExistingEnvironment, "!cmd", "fin"},
	}

	for _, tc := range testCases {
		config.Conf = config.NewMockConfig(tc.option)
		parser = command.NewMockParser()

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
		{NoAccess, "!cmd", "Access denied"},
		{FullAccess, "!cmd", "fin"},
	}

	for _, tc := range testCases {
		config.Conf = config.NewMockConfig(tc.option)
		parser = command.NewMockParser()

		runTest(t, tc.incomingText, tc.outgoingText)
	}
}

func TestGetUser(t *testing.T) {
	users = []slack.User{
		slack.User{ID: "u1", Name: "user1"},
		slack.User{ID: "u2", Name: "user2"},
	}

	testCases := []struct {
		id       string
		expected string
	}{
		{"u1", "user1"},
		{"u2", "user2"},
		{"unknown", ""},
	}

	for _, tc := range testCases {
		assert.Equal(t, GetUser(tc.id).Name, tc.expected)
	}
}

func TestGetChannel(t *testing.T) {

	channels = []slack.Channel{
		makeChannel("ch1", "channel1"),
		makeChannel("ch2", "channel2"),
	}

	testCases := []struct {
		id       string
		expected string
	}{
		{"ch1", "channel1"},
		{"ch2", "channel2"},
		{"unknown", ""},
	}

	for _, tc := range testCases {
		assert.Equal(t, GetChannel(tc.id).Name, tc.expected)
	}
}

func runTest(t *testing.T, incomingText string, outgoingText string) {
	if outgoingText == "" {
		return
	}

	rtm.IncomingEvents <- newMessageEvent(incomingText)
	rsp, err := server.waitForRequest()
	if err != nil {
		t.Fatal(err)
	}
	msg := getOutgoingMessage(rsp)

	assert.Equal(t, outgoingText, msg.Text)
}

func makeChannel(id string, name string) slack.Channel {
	ch := slack.Channel{}
	ch.ID = id
	ch.Name = name

	return ch
}
