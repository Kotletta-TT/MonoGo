include .env

server:
	GOOS=linux GOARCH=amd64 go build -o cmd/server/server cmd/server/main.go

agent:
	GOOS=linux GOARCH=amd64 go build -o cmd/agent/agent cmd/agent/main.go

build:
	make server
	make agent

run_serv:
	ADDRESS="${ADDRESS}" LOG_LEVEL="${LOG_LEVEL}" LOG_PATH="${LOG_PATH}" go run cmd/server/main.go

run_agent:
	ADDRESS="${ADDRESS}" REPORT_INTERVAL="${REPORT_INTERVAL}" POLL_INTERVAL="${POLL_INTERVAL}" go run cmd/agent/main.go

.PHONY: build server agent