package slack

import (
	"errors"
	"fmt"
	"github.com/apex/log"

	"kubot/command"
	"kubot/config"

	"github.com/nlopes/slack"
)

var (
	startOptions []slack.Option
	rtm          *slack.RTM
	parser       command.CommandParser
	users        []slack.User
	channels     []slack.Channel
)

func init() {
	startOptions = []slack.Option{}
	token := config.Conf.GetSlackToken()
	api := slack.New(token, startOptions...)
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
	log.WithField("type", fmt.Sprintf("%T", e.Data)).Info("incoming event")

	switch ev := e.Data.(type) {
	case *slack.MessageEvent:
		cmd, err := parser.Parse(ev.Text)
		if err != nil {
			return // Fail silently
		}

		env, err := config.Conf.GetEnvironmentByChannel(getChannel((ev.Channel)).Name)
		if err != nil {
			handleError(err, ev.Channel)
			return
		}

		if !config.Conf.HasAccess(getUser(ev.User).Profile.Email, env.Name) {
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
