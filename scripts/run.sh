#!/bin/bash

set -e


GOARCH=amd64 GOOS=windows go build -o build/app.exe ./cmd/sven  && powershell.exe -Command ".\\build\\app.exe $@"
rm -f build/app.exe
rm -f audio.mp3
