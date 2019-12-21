package main

import (
	"kubot/api"
	"kubot/slack"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var apiPort string

func init() {
	rootCmd = &cobra.Command{
		Use: "kubot",
		Run: func(c *cobra.Command, args []string) { run() },
	}
	rootCmd.Flags().StringVarP(&apiPort, "apiPort", "p", "", "enable an api web server to accept requests")
}

func main() {
	rootCmd.Execute()
}

func run() {
	log.SetFormatter(&log.JSONFormatter{})

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go api.Start(apiPort)
	go slack.Start()
	<-stop
}
