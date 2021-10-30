# HTTP to WS proxy

Minimal proxy available as package or as binary.

## Usecase

I quickly created it to solve my problem - I needed to connect to Binance BSC blockchain testnet node through websocket (because app I was contributing to supported only websocket providers).
But testnet node did not supported websocket - so I created this little script.

And I am sure I am recreating a wheel here - probably it was already implemented many times by other projects...

## Package Usage

`cmd/main.go` and `proxy_test.go` also consume code as package so you can refer to those as examples as well.

```go
p := http2ws.Proxy{
  TargetOpts: http2ws.TargetOpts{
    Method: "POST",
    URL:    "http://some-http-serv",
    Headers: map[string]string{
      "Content-Type": "application/json",
    },
  },
  Log: log.Default(),
  WSUpgrader: websocket.Upgrader{
    // allow all origins
    CheckOrigin: func(r *http.Request) bool { return true },
  },
  HTTPClient: http.Client{
    Timeout: time.Second,
  },
}

http.ListenAndServe(args.addr, &p)
```

## Executable Usage

Print help

```
http2ws --help
```

Forward incoming websocket messages to `target-url` as HTTP `POST` requests with adding header `Content-Type:application/json` (multiple headers can be specified, see help)

```
http2ws -addr=:4000 -target-url="http://some-http-serv" -target-method=POST -target-headers="Content-Type:application/json"
```
