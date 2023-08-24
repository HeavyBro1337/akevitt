#!/bin/bash
env GOARCH="amd64" GOOS="linux" go build -o akevitt-generator-linux-amd64
env GOARCH="386" GOOS="windows" go build -o akevitt-generator-win-x86.exe
