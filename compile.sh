#!/bin/bash

echo "Linux AMD64"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o files/historian_linux_amd64
echo "Linux i386"
CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o files/historian_linux_386
echo "Linux ARM5"
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -o files/historian_linux_arm5
echo "Linux ARM6"
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -o files/historian_linux_arm6
echo "Linux ARM7"
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o files/historian_linux_arm7
echo "Darwin i386"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o files/historian_darwin_386
echo "Darwin AMD 64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o files/historian_darwin_amd64