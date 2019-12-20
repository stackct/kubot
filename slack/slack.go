package slack

import (
	"kubot/command"
	"log"
	"os"

	"github.com/nlopes/slack"
)

var (
	startOptions []slack.Option
	rtm          *slack.RTM
	parser       command.SlackCommandParser
)

func init() {
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

	for e := range rtm.IncomingEvents {
		handleEvent(e)
	}
}

func handleEvent(e slack.RTMEvent) {
	switch ev := e.Data.(type) {
	case *slack.MessageEvent:
		cmd, err := parser.Parse(ev.Text)
		if err != nil {
			rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Channel))
			return
		}

		out := make(chan string)
		go cmd.Execute(out)

		for msg := range out {
			rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
		}
	}
}
