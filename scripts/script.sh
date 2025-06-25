#!/bin/bash


GOARCH=amd64 GOOS=windows go build -o build/app.exe . && powershell.exe -Command ".\\build\\app.exe $@"
rm build/app.exe
rm adam.mp3
