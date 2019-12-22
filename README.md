# kubot
Kubot is a Slack integration (bot) for executing deployments.

[![Build Status](https://travis-ci.org/dotariel/kubot.svg?branch=master)](https://travis-ci.org/dotariel/kubot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotariel/kubot)](https://goreportcard.com/report/github.com/dotariel/kubot)
[![Code Coverage](https://codecov.io/gh/dotariel/kubot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotariel/kubot)
[![Docker Build](https://img.shields.io/docker/cloud/automated/dotariel/kubot)](https://hub.docker.com/r/dotariel/kubot)
[![Docker Pulls](https://img.shields.io/docker/pulls/dotariel/kubot.svg)](https://hub.docker.com/r/dotariel/kubot)

## Commands
Kubot's configurable command system allows commands to be declared via configuration.

Example:
```
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

## Docker

```
docker run -v config.yml:/config.yml -it dotariel/kubot
```