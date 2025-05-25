#! /usr/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

source $SCRIPT_DIR/set-env.sh

sudo docker run \
  --name $DOCKER_NAME -e POSTGRES_PASSWORD=$DB_PASSWORD \
  -p $PORT:$PORT \
  -d postgres

# To remove the docker image:
# sudo docker rm -f some-db
