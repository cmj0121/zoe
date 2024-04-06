#! /usr/bin/env sh
set -e

mkdir -p /svc/zoe

zoe -vv migrate sqlite3:///svc/zoe/records.db
zoe -vv launch --config zoe.yaml
