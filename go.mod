module kubot

// Temporarily override until PR is applied (https://github.com/nlopes/slack/pull/636)
replace github.com/nlopes/slack v0.6.0 => github.com/acaloiaro/slack v0.6.3-0.20191210002151-2cc5dc1c8f87

go 1.13

require (
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.2.0
	github.com/nlopes/slack v0.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.2.2
	gopkg.in/yaml.v2 v2.2.7
)
