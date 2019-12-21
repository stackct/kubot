package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kubot/command"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
		log.Error("Failed to parse cmd request; reason=", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := make(chan string)
	go cmd.Execute(out)
	var messages []string
	for msg := range out {
		log.Info(msg)
		messages = append(messages, msg)
	}
	json.NewEncoder(w).Encode(messages)
}
