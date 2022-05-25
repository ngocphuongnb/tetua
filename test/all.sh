#!/bin/bash

# go test ./...
go test -coverprofile $PWD/test/coverage.txt ./... && go tool cover -html=$PWD/test/coverage.txt -o $PWD/test/coverage.html
# echo $?