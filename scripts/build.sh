#!/usr/bin/env bash
set -exuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc

GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/varify ./src/nginx/varify
GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/supply ./src/nginx/supply/cli
GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/finalize ./src/nginx/finalize/cli
