package http2ws_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ulexxander/go-http2ws"
)

func TestProxy(t *testing.T) {
	type message struct {
		Hello string
	}

	msg := message{Hello: "there"}

	targetChan := make(chan message)

	targetServ := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body message
		err := json.NewDecoder(r.Body).Decode(&body)
		assertNoError(t, err, "decoding target json")
		targetChan <- body
		w.WriteHeader(http.StatusOK)
	}))

	p := http2ws.Proxy{
		TargetOpts: http2ws.TargetOpts{
			Method: "GET",
			URL:    targetServ.URL,
		},
	}
	proxyServ := httptest.NewServer(&p)
	defer proxyServ.Close()

	conn, res, err := websocket.DefaultDialer.Dial(httpToWs(proxyServ.URL), nil)
	assertNoError(t, err, "dialing")
	assertStatusCode(t, res, 101)
	defer conn.Close()

	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("message %d", i), func(t *testing.T) {
			err := conn.WriteJSON(msg)
			assertNoError(t, err, "writing json payload")

			select {
			case result := <-targetChan:
				if result.Hello != msg.Hello {
					t.Fatalf("expected result.Hello to be %s, got: %s", msg.Hello, result.Hello)
				}
			case <-time.After(time.Second):
				t.Fatalf("not proxied in one second")
			}
		})
	}
}

func httpToWs(httpUrl string) string {
	return strings.Replace(httpUrl, "http", "ws", 1)
}

func assertStatusCode(t *testing.T, r *http.Response, exp int) {
	if r.StatusCode != exp {
		t.Fatalf("expected status code %d, got: %d", exp, r.StatusCode)
	}
}

func assertNoError(t *testing.T, err error, desc string) {
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", desc, err)
	}
}
