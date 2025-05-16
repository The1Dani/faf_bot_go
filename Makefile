all: main

main: ./cmd/bot/*
	go build ./cmd/bot

run: main
	scripts/db-start.sh && source scripts/set-env.sh && ./bot

stop-db:
	scripts/db-stop.sh
