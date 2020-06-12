package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kubot/command"
	"kubot/config"
	"net/http"

	"github.com/apex/log"

	"github.com/gorilla/mux"
)

func Start(port string) {
	if "" != port {
		r := mux.NewRouter()
		r.HandleFunc("/", Execute).Methods("POST")
		http.Handle("/", r)
		log.Info(fmt.Sprintf("Now listening on: %s", port))
		http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	}
}

func Execute(w http.ResponseWriter, r *http.Request) {
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	cmd, err := command.NewSlackCommandParser().Parse(string(bodyBuffer))
	if err != nil {
		log.WithField("reason", err.Error()).Error("Failed to parse cmd request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	env, err := config.Conf.GetEnvironmentByChannel("api")
	if err != nil {
		log.WithField("reason", err.Error()).Error("Failed to parse environment config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := make(chan string)
	go cmd.Execute(out, command.Context{Environment: *env})
	var messages []string
	for msg := range out {
		log.Info(msg)
		messages = append(messages, msg)
	}
	json.NewEncoder(w).Encode(messages)
}
