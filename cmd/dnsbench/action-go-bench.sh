#!/bin/sh
go tool vet > /dev/null 2>&1
export GOGC=off
GOGC="off" go test -bench . -mod=vendor
if [ -r main.go ]; then go run main.go; fi
