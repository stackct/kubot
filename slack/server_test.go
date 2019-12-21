package slack

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

	log.Info("mock server started", zap.String("httpUrl", server.httpURL()), zap.String("wsUrl", server.wsURL()))

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
	log.Debug("handling request", zap.String("url", r.URL.String()))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade failed", zap.Error(err))
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error("read failed", zap.Error(err))
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
