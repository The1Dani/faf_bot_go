#! /usr/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd $SCRIPT_DIR

export ../env.sh
sudo TOKEN=$TELEGRAM_APITOKEN docker compose up -d
