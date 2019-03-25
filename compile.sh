#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o files/historian_linux_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o files/historian_linux_386
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o files/historian_linux_arm5
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o files/historian_linux_arm6
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o files/historian_linux_arm7
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o files/historian_darwin_386
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o files/historian_darwin_amd64