module kubot

// Temporarily override until PR is applied (https://github.com/nlopes/slack/pull/636)
replace github.com/nlopes/slack v0.6.0 => github.com/acaloiaro/slack v0.6.3-0.20191210002151-2cc5dc1c8f87

go 1.15

require (
	github.com/apex/log v1.9.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/nlopes/slack v0.6.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)
