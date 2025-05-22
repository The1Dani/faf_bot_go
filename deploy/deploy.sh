#! /usr/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

source ../env.sh
sudo docker compose down
sudo docker compose build
sudo TOKEN=$TELEGRAM_APITOKEN docker compose up -d
