#! /usr/bin/env sh
set -e

mkdir -p /svc/zoe

zoe -vv migrate --database sqlite3:///svc/zoe/records.db --folder assets/migrations
zoe -vv launch --config zoe.yaml
