package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type TargetOpts struct {
	URL     string
	Method  string
	Headers map[string]string
}

type Proxy struct {
	TargetOpts TargetOpts
	Log        *log.Logger
	Upgrader   websocket.Upgrader
	HTTPClient http.Client
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := p.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.Printf("failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	for {
		mt, payload, err := conn.ReadMessage()
		if err != nil {
			p.Printf("failed to read message: %v", err)
			return
		}
		if mt != websocket.TextMessage {
			p.Printf("not a text message (1): %d", mt)
			continue
		}

		p.Println("incoming request, len:", len(payload))

		response, err := p.request(payload)
		if err != nil {
			response = []byte(err.Error())
		}
		if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
			p.Printf("failed to write message: %s", err)
			return
		}
	}
}

func (p *Proxy) request(payload []byte) ([]byte, error) {
	reqBody := bytes.NewReader(payload)
	r, err := http.NewRequest(p.TargetOpts.Method, p.TargetOpts.URL, reqBody)
	for k, v := range p.TargetOpts.Headers {
		r.Header.Add(k, v)
	}
	if err != nil {
		return nil, fmt.Errorf("initializing new request: %v", err)
	}
	res, err := p.HTTPClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("sending request: %v", err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %v", err)
	}
	return resBody, nil
}

func (p *Proxy) Printf(format string, v ...interface{}) {
	if p.Log != nil {
		p.Log.Printf(format, v...)
	}
}

func (p *Proxy) Println(v ...interface{}) {
	if p.Log != nil {
		p.Log.Println(v...)
	}
}
