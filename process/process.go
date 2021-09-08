package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/apex/log"
)

func Start(name string, args []string, replacementArgs map[string]string, environmentVariables map[string]string) ([]byte, error) {
	resolvedArgs := Interpolate(args, replacementArgs)
	log.WithField("name", name).WithField("args", resolvedArgs).Info("executing command")
	cmd := exec.Command(name, resolvedArgs...)
	cmd.Env = os.Environ()
	for k, v := range environmentVariables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	out, err := cmd.CombinedOutput()

	logEntry := log.WithField("name", name).WithField("args", resolvedArgs).WithField("output", string(out))
	if nil != err {
		logEntry.WithError(err).WithField("stackTrace", string(debug.Stack())).Error("command failed")
	} else {
		logEntry.Info("command completed")
	}
	return out, err
}

func Interpolate(args []string, replacementArgs map[string]string) []string {
	result := strings.Join(args, ",")
	for k, v := range replacementArgs {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", k), v)
	}
	return strings.Split(result, ",")
}
