# kubot
Kubot is a Slack integration (bot) for executing deployments.

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
- !kick <pod>
- restart a pod

- !restart <product>
- kubectl -n <product> rollout restart deployment/<product>

- !secret <product>
- helm apply secret

- !deploy <product> <version>
- helm repo update
- checkout the environment repo
- update and push the product version change
- helm upgrade product repo/chartname --wait --timeout 900 --version version -f environment/product/yaml
- Log all command output to kibana and scrub any sensitive data
- Upon completion, post a message to slack