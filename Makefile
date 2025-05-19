all: main

main: ./cmd/bot/*
	go build -o ./build/bot ./cmd/bot 

run: main
	scripts/db-start.sh && source scripts/set-env.sh && ./build/bot

stop-db:
	scripts/db-stop.sh
