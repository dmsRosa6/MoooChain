#!/bin/bash

set -e

CONTAINER_NAME="moochain-redis"
SCRIPTS_FILE_NAME="utilslib.lua"

echo "🧹 Cleaning up any existing container..."
docker rm -f $CONTAINER_NAME 2>/dev/null || true


echo "📚 Starting Redis container..."
docker run --name $CONTAINER_NAME \
  -p 6379:6379 \
  -d redis

echo "⏳ Waiting..."
sleep 5

echo "📄 Loading Scripts..."
redis-cli -x function load < $SCRIPTS_FILE_NAME