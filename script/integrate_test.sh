#!/usr/bin/env bash

docker compose -f script/cache_test_compose.yml down
docker compose -f script/cache_test_compose.yml up -d
go test -race ./... -tags=e2e
docker compose -f script/cache_test_compose.yml down
