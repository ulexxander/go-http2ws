FROM golang:1.17-alpine AS build
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/http2ws/main.go

FROM alpine
WORKDIR /http2ws
COPY --from=build /build/main .
ENTRYPOINT ["./main"]
