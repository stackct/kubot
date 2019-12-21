package main

import (
	"kubot/api"
	"kubot/config"
	"kubot/slack"
	"os"
	"os/signal"
	"syscall"

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
	defer config.Log.Sync()
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go api.Start(apiPort)
	go slack.Start()
	<-stop

}
