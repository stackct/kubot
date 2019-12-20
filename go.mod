module kubot

require github.com/dotariel/kubot/slack v0.0.0

// Temporarily override until PR is applied (https://github.com/nlopes/slack/pull/636)
replace github.com/nlopes/slack v0.6.0 => github.com/acaloiaro/slack v0.6.3-0.20191210002151-2cc5dc1c8f87

replace github.com/dotariel/kubot/slack => ./slack

replace kubot/command => ./command

go 1.13
