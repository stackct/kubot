# kubot
Kubot is a Slack integration (bot) for executing deployments.

[![Build Status](https://travis-ci.org/dotariel/kubot.svg?branch=master)](https://travis-ci.org/dotariel/kubot)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotariel/kubot)](https://goreportcard.com/report/github.com/dotariel/kubot)
[![Code Coverage](https://codecov.io/gh/dotariel/kubot/branch/master/graph/badge.svg)](https://codecov.io/gh/dotariel/kubot)
[![Docker Build](https://img.shields.io/docker/cloud/automated/dotariel/kubot)](https://hub.docker.com/r/dotariel/kubot)
[![Docker Pulls](https://img.shields.io/docker/pulls/dotariel/kubot.svg)](https://hub.docker.com/r/dotariel/kubot)

# Local developer setup
export KUBOT_SLACK_TOKEN=secrettoken

# TODO
- Create a docker image to run kubot
- Inject secrets to pull helm charts, an environment specific slack token, kubit environment, vault secrets
- Create helm charts for deploying kubot
- Update kubit toxic job to auto deploy kubot
- Before running any command, perform an authorization check by verifying the slack user exists in a config file
- Add a !release <product> command to toxic to perform a make release to increment the version and tag the product

# kubot commands
- !deploy <product> <version>

- !kick <pod>
- restart a pod

- !restart <product>
- kubectl -n <product> rollout restart deployment/<product>

## Local Setup
Kubot configuration is managed through a configuration file located at `KUBOT_CONFIG` (see [sample file](config/resources/kubot.yml)).

## Docker

```
docker run -e KUBOT_SLACK_TOKEN=xxxxx -v config.yml:/config.yml -it dotariel/kubot
```