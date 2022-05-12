#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Illegal number of parameters"
fi

rootDir=$PWD

cd $1
go test -coverprofile $rootDir/test/coverage.txt
go tool cover -html=$rootDir/test/coverage.txt -o $rootDir/test/coverage.html
cd $rootDir