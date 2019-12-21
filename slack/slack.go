package slack

import (
	"errors"
	"kubot/command"
	"log"
	"os"

	"github.com/nlopes/slack"
	"kubot/config"
)

var (
	Conf         config.Configurator
	startOptions []slack.Option
	rtm          *slack.RTM
	parser       command.CommandParser
	users        []slack.User
	channels     []slack.Channel
)

func init() {
	Conf, _ = config.ParseFile(os.Getenv("KUBOT_CONFIG"))

	startOptions = []slack.Option{
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	}

	api := slack.New(os.Getenv("KUBOT_SLACK_TOKEN"), startOptions...)
	rtm = api.NewRTM()
	parser = command.NewSlackCommandParser()
}

func Start() {
	go rtm.ManageConnection()

	users, _ = rtm.GetUsers()
	channels, _ = rtm.GetChannels(true)

	for e := range rtm.IncomingEvents {
		handleEvent(e)
	}
}

func handleEvent(e slack.RTMEvent) {
	switch ev := e.Data.(type) {
	case *slack.MessageEvent:
		cmd, err := parser.Parse(ev.Text)
		if err != nil {
			handleError(err, ev.Channel)
			return
		}

		env, err := Conf.GetEnvironmentByChannel(getChannel((ev.Channel)).Name)
		if err != nil {
			handleError(err, ev.Channel)
			return
		}

		if !Conf.HasAccess(getUser(ev.User).Profile.Email, env.Name) {
			handleError(errors.New("Authorization denied"), ev.Channel)
			return
		}

		out := make(chan string)
		go cmd.Execute(out)

		for msg := range out {
			rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
		}
	}
}

func handleError(err error, channel string) {
	rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), channel))
}

func getUser(id string) *slack.User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}

	return &slack.User{}
}

func getChannel(id string) slack.Channel {
	for _, channel := range channels {
		if channel.ID == id {
			return channel
		}
	}

	return slack.Channel{}
}
