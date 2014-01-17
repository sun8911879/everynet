#!/usr/bin/env bash

chmod +x $0

File="client"

cd client

if [ ! -f "$File" ]; then
go build client.go
fi

./client
