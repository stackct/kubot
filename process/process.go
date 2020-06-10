package process

import (
	"fmt"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/apex/log"
)

func Start(name string, args []string, replacementArgs map[string]string) error {
	resolvedArgs := Interpolate(args, replacementArgs)
	log.WithField("name", name).WithField("args", resolvedArgs).Info("executing command")
	out, err := exec.Command(name, resolvedArgs...).Output()

	logEntry := log.WithField("name", name).WithField("args", resolvedArgs).WithField("stdout", string(out))
	if nil != err {
		logEntry.WithError(err).WithField("stackTrace", string(debug.Stack())).Error("command failed")
	} else {
		logEntry.Info("command completed")
	}
	return err
}

func Interpolate(args []string, replacementArgs map[string]string) []string {
	result := strings.Join(args, ",")
	for k, v := range replacementArgs {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", k), v)
	}
	return strings.Split(result, ",")
}
