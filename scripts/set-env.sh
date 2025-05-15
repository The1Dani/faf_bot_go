#! /usr/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

export PORT=5432
export HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=11111
export DOCKER_NAME=some-db

source $SCRIPT_DIR/.env
export FAFBOT_PGSQL_CONNECTION="postgresql://$HOST:$PORT?user=$DB_USER&password=$DB_PASSWORD&sslmode=disable"
export TELEGRAM_APITOKEN="$TELEGRAM_APITOKEN"