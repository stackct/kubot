package command

import (
	"fmt"
	"github.com/apex/log"
	"io"
	"kubot/config"
	"os/exec"
	"strings"
)

type Command interface {
	Execute(output chan string)
}

type StdoutWriter struct{}
type StderrWriter struct{}

var CommandStdoutWriter io.Writer
var CommandStderrWriter io.Writer

func init() {
	CommandStdoutWriter = &StdoutWriter{}
	CommandStderrWriter = &StderrWriter{}
}

func (w StdoutWriter) Write(msg []byte) (n int, err error) {
	log.Info(string(msg))
	return len(msg), nil
}

func (w StderrWriter) Write(msg []byte) (n int, err error) {
	log.Error(string(msg))
	return len(msg), nil
}

func Execute(name string, configOverrides map[string]string) error {
	command, err := config.Conf.GetCommand(name)
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
		cmd := exec.Command(command.Commands[i].Name, Interpolate(command.Commands[i], commandConfig)...)
		cmd.Stdout = CommandStdoutWriter
		cmd.Stderr = CommandStderrWriter

		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func Interpolate(command config.Command, replacementArgs map[string]string) []string {
	result := strings.Join(command.Args, ",")
	for k, v := range replacementArgs {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", k), v)
	}
	return strings.Split(result, ",")
}
