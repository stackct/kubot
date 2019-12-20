package main

import (
	"github.com/dotariel/kubot/slack"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go slack.Start()
	<-stop
}
