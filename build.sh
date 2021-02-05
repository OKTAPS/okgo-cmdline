#!/bin/bash
go mod tidy
env GOOS=linux GOARCH=amd64 go build -o ./build/linux-amd64/okgo
env GOOS=linux GOARCH=386 go build -o ./build/linux-386/okgo
env GOOS=linux GOARCH=arm go build -o ./build/linux-arm/okgo
env GOOS=windows GOARCH=386 go build -o ./build/linux-386/okgo
env GOOS=windows GOARCH=amd64 go build -o ./build/windows-amd64/okgo
env GOOS=darwin GOARCH=amd64 go build -o ./build/darwin-amd64/okgo
env GOOS=darwin GOARCH=386 go build -o ./build/darwin-386/okgo
