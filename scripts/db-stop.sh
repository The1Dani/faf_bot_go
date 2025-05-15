#! /usr/bin/bash
    
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

source SCRIPT_DIR/set-env.sh

sudo docker stop $DOCKER_NAME