#! /usr/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

source $SCRIPT_DIR/set-env.sh

pg_isready --host=$HOST --port=$PORT --username=$DB_USER

if [[ $? -eq 0 ]]; then
  echo "[SCRIPT] DB is already started"
  exit 0
else
  echo "[SCRIPT] Starting $DOCKER_NAME"
  sudo docker start $DOCKER_NAME
fi

if [[ $? -eq 0 ]]; then
  echo "[SCRIPT] DB already exists starting"
  exit 0
else
  echo "[SCRIPT] creating db instance named $DOCKER_NAME"
  bash $SCRIPT_DIR/db-first-start.sh
fi
