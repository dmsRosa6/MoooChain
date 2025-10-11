#!/bin/bash

set -e

CONTAINER_NAME="moochain-redis"
SCRIPTS_FILE_NAME="utilslib.lua"

echo "Loading Scripts..."
redis-cli -x function load replace < $SCRIPTS_FILE_NAME