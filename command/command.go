package command

import (
	"kubot/config"
	"kubot/process"
)

type Command interface {
	Execute(output chan string, context Context)
}

func Execute(name string, product string, configOverrides map[string]string, environmentVariables map[string]string) error {
	command, err := config.Conf.GetCommand(name, product)
	if err != nil {
		return err
	}

	commandConfig := config.Conf.GetCommandConfig()
	for k, v := range command.Config {
		commandConfig[k] = v
	}
	for k, v := range configOverrides {
		commandConfig[k] = v
	}

	for i := 0; i < len(command.Commands); i++ {
		if err := process.Start(command.Commands[i].Name, command.Commands[i].Args, commandConfig, environmentVariables); err != nil {
			return err
		}
	}

	return nil
}
