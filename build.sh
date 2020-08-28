#! /bin/bash
# Copyright 2018 Kuei-chun Chen. All rights reserved.

DEP=`which dep`
if [ "$DEP" == "" ]; then
    echo "dep command not found"
    exit
fi

if [ -d vendor ]; then
    UPDATE="-update"
fi

$DEP ensure $UPDATE

export version="0.2.0"
mkdir -p build
# env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$version" -o build/argos-linux-x64 argos.go
env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$version" -o build/argos-osx-x64 argos.go
# env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$version" -o build/argos-win-x64.exe argos.go
