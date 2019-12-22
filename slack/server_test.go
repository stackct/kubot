package slack

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"time"
)

type MockServer struct {
	*httptest.Server
	requests  chan string
	timeoutMs int
}

func newMockServer() *MockServer {
	server := &MockServer{
		Server:    httptest.NewServer(nil),
		requests:  make(chan string),
		timeoutMs: 10,
	}

	http.HandleFunc("/rtm", server.websocket)
	http.HandleFunc("/rtm.connect", server.connect)

	log.WithField("httpUrl", server.httpURL()).WithField("wsUrl", server.wsURL()).Info("mock server started")

	return server
}

func (ms *MockServer) httpURL() string {
	return "http://" + ms.Listener.Addr().String() + "/"
}

func (ms *MockServer) wsURL() string {
	return "ws://" + ms.Listener.Addr().String() + "/rtm"
}

func (ms *MockServer) record(message string) {
	ms.requests <- message
}

func (ms *MockServer) waitForRequest() (string, error) {
	select {
	case msg := <-ms.requests:
		return msg, nil
	case <-time.After(time.Duration(ms.timeoutMs) * time.Millisecond):
		return "", errors.New("Timed out waiting for request")
	}
}

func (ms *MockServer) connect(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf(`{ "ok": true, "url":"%v"}`, ms.wsURL())
	bytes := []byte(response)
	w.Write(bytes)
}

func (ms *MockServer) websocket(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL).Debug("handling request")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("reason", err).Error("upgrade failed")
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.WithField("reason", err).Error("read failed")
			break
		}

		// Right here is where we can respond
		ms.record(string(message))

		// TODO: Write back
	}
}

func (ms *MockServer) stop() {
	ms.Close()
}
