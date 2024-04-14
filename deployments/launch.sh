#! /usr/bin/env sh
set -e

mkdir -p /tmp/zoe

zoe -vv migrate sqlite3:///tmp/zoe/records.db
zoe -vv launch --config zoe.yaml
