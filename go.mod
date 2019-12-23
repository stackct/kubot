module kubot

// Temporarily override until PR is applied (https://github.com/nlopes/slack/pull/636)
replace github.com/nlopes/slack v0.6.0 => github.com/acaloiaro/slack v0.6.3-0.20191210002151-2cc5dc1c8f87

go 1.13

require (
	github.com/acaloiaro/slack v0.6.3 // indirect
	github.com/apex/log v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.2.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.7
)
