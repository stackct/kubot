package process

import (
	"fmt"
	"github.com/apex/log"
	"io"
	"os/exec"
	"strings"
)

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

func Start(name string, args []string, replacementArgs map[string]string) error {
	resolvedArgs := Interpolate(args, replacementArgs)
	cmd := exec.Command(name, resolvedArgs...)
	cmd.Stdout = CommandStdoutWriter
	cmd.Stderr = CommandStderrWriter

	log.WithField("name", name).WithField("args", resolvedArgs).Info("executing command")

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func Interpolate(args []string, replacementArgs map[string]string) []string {
	result := strings.Join(args, ",")
	for k, v := range replacementArgs {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", k), v)
	}
	return strings.Split(result, ",")
}
