package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/ulexxander/go-http2ws"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %s\n", err)
	}
}

type args struct {
	addr          string
	targetURL     string
	targetMethod  string
	targetHeaders string
}

func parseArgs() (*args, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	addr := flagSet.String("addr", ":80", "address to listen on")
	targetURL := flagSet.String("target-url", "", "target url of proxy")
	targetMethod := flagSet.String("target-method", "GET", "http method of request that is sent to target")
	targetHeaders := flagSet.String("target-headers", "", "list of '"+headersSeparator+"' separated headers to attach to target request")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			flagSet.Usage()
		}
		return nil, err
	}

	if *addr == "" {
		return nil, errors.New("addr cannot be empty")
	}
	if *targetURL == "" {
		return nil, errors.New("target-url cannot be empty")
	}
	if *targetMethod == "" {
		return nil, errors.New("target-method cannot be empty")
	}

	return &args{
		addr:          *addr,
		targetURL:     *targetURL,
		targetMethod:  *targetMethod,
		targetHeaders: *targetHeaders,
	}, nil
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return fmt.Errorf("parsing args: %v", err)
	}

	headers, err := parseHeaders(args.targetHeaders)
	if err != nil {
		return fmt.Errorf("parsing headers: %v", err)
	}

	p := http2ws.Proxy{
		TargetOpts: http2ws.TargetOpts{
			Method:  args.targetMethod,
			URL:     args.targetURL,
			Headers: headers,
		},
		Log: log.Default(),
		WSUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	log.Println("listening on", args.addr)
	if err := http.ListenAndServe(args.addr, &p); err != nil {
		return fmt.Errorf("failed listening: %v", err)
	}
	return nil
}

const headersSeparator = "|"

func parseHeaders(arg string) (map[string]string, error) {
	parts := strings.Split(arg, "|")
	headers := map[string]string{}
	for _, p := range parts {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid value: %s", p)
		}
		headers[kv[0]] = kv[1]
	}
	return headers, nil
}
