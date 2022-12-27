package slack

import (
	"kubot/command"
	"kubot/config"
	"net/http"
	"testing"

	"github.com/slack-go/slack"
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
		{ID: "_user", Profile: slack.UserProfile{Email: "foo@invalid.co"}},
	}
	startOptions = append(startOptions, slack.OptionAPIURL(server.httpURL()))
	rtm = slack.New("", startOptions...).NewRTM()

	go Start()

	return Teardown
}

func Teardown() {
	parser = command.NewSlackCommandParser()
}

func TestStart_Parse_Command_Error(t *testing.T) {
	defer BeforeEach(t)()

	testCases := []struct {
		option       command.MockParserOption
		incomingText string
		outgoingText string
	}{
		{ParseUnknownCommandError, "any", ""},
		{ParseCommandArgumentError, "any", "command argument error"},
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

func TestStart_Loads_Channels(t *testing.T) {
	defer BeforeEach(t)()

	parser = command.NewMockParser()

	http.HandleFunc("/conversations.list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "ok": true, "channels": [ { "id": "123FOO", "name": "foo" } ] }`))
	})

	runTest(t, "!any", "Ready")

	rsp, err := server.waitForRequest()
	if err != nil {
		t.Fatal(err)
	}
	msg := getOutgoingMessage(rsp)
	assert.Equal(t, "fin", msg.Text)

	assert.Equal(t, 1, len(channels))
	assert.Equal(t, "foo", channels[0].Name)
}

func TestGetUser(t *testing.T) {
	users = []slack.User{
		{ID: "u1", Name: "user1"},
		{ID: "u2", Name: "user2"},
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
