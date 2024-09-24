#! /usr/bin/env sh
set -e

# replace the config file with env vars
envsubst < /app/zoe.yml > zoe.yml
mkdir -p .ssh/

/usr/local/bin/zoe -v -c zoe.yml ssh
