module github.com/dotariel/kubot/slack

go 1.13

require (
	github.com/gorilla/websocket v1.4.1
	github.com/nlopes/slack v0.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	kubot/command v0.0.0
)

replace kubot/command => ../command
