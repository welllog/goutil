#!/bin/sh

protoc -I./ -I${GOPATH}/src \
		--go_out=./ \
    --go_opt=paths=source_relative \
    ./hello.proto