#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

VER=0.1.0beta
ID=$(git rev-parse HEAD | cut -c1-7)

if [[ $1 = "-mac" ]]; then
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv gcollector bin/gcollector.mac

    cd cmd/benchmark
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv benchmark ../../bin/benchmark.mac
else
    go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    #go build -race -v -ldflags "-X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv gcollector bin/gcollector.linux
    bin/gcollector.linux -v

    cd cmd/benchmark
    go build -ldflags "-X github.com/nicholaskh/golib/server.VERSION $VER -X github.com/nicholaskh/golib/server.BuildID $ID -w"
    mv benchmark ../../bin/benchmark.linux
fi
