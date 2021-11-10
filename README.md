# ![kubot][logo] kubot

Kubot is a Slack integration (bot) for executing deployments.

[![Build Status](https://travis-ci.org/dotariel/kubot.svg?branch=master)](https://travis-ci.org/dotariel/kubot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotariel/kubot)](https://goreportcard.com/report/github.com/dotariel/kubot)
[![Code Coverage](https://codecov.io/gh/dotariel/kubot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotariel/kubot)

## Commands

Kubot's configurable command system allows commands to be declared via configuration.

Example:

```
environments:
  - name: local
    channel: env_local
    users:
      - foo@bar.com
    variables:
      - foo: "bar"

commandConfig:
  dir: /tmp

commands:
  - name: list
    commands:
      - name: "ls"
        args:
          - "${dir}"
```

## Local Setup

Kubot configuration is managed through a configuration file located at `KUBOT_CONFIG` (see [sample file](config/resources/kubot.yml)).

## API Setup

Kubot can be configured to listen for slack commands on a local port for easier developer testing.

```
go build && ./kubot -p 8080
curl -X POST --data ".deploy foo 1.0.0" http://localhost:8080
```

## Docker

```
docker run -v config.yml:/config.yml -it dotariel/kubot
```

[logo]: assets/kubot-24x24.png "kubot"
