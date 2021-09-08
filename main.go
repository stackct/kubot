package main

import (
	"kubot/api"
	"kubot/config"
	"kubot/process"
	"kubot/slack"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"

	"github.com/spf13/cobra"
)

var (
	version   string
	buildDate string
	rootCmd   *cobra.Command
	apiPort   string
	logFile   *os.File
)

func init() {
	rootCmd = &cobra.Command{
		Use: "kubot",
		Run: func(c *cobra.Command, args []string) { run() },
	}
	rootCmd.Flags().StringVarP(&apiPort, "apiPort", "p", "", "enable an api web server to accept requests")
}

func main() {
	logging := config.Conf.GetLogging()
	logFile, _ = config.InitLogging(logging.File, logging.Level)
	defer logFile.Close()

	rootCmd.Execute()
}

func run() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for _, init := range config.Conf.GetInit() {
		log.
			WithField("command", init.Name).
			WithField("args", init.Args).
			Info("executing init command")

		if _, err := process.Start(init.Name, init.Args, config.Conf.GetCommandConfig(), map[string]string{}); err != nil {
			log.
				WithField("command", init.Name).
				WithField("args", init.Args).
				WithError(err).
				Error("init command failed")
		}
	}

	go api.Start(apiPort)
	go slack.Start()

	log.WithField("version", version).WithField("buildDate", buildDate).Info("Application Started")

	<-stop
}
