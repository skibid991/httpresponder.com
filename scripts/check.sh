#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
node --test sdk/javascript/client.test.mjs
PYTHONPATH=sdk/python python3 -m unittest discover -s sdk/python -p 'test_*.py'
go -C sdk/go test ./...
go -C examples/go build -o "${TMPDIR:-/tmp}/httpresponder-example" .
