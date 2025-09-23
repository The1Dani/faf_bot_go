#! /usr/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

source ../env.sh
docker compose down
docker compose build
TOKEN=$TELEGRAM_APITOKEN docker compose up -d
