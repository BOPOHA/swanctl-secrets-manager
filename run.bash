#!/bin/bash
set -eu
go mod tidy && go mod vendor
gofmt -w *.go

BIN_NAME=./bin/go-secrets-create
CORP_DOMAIN=vorona.me
SUB_DOMAIN=vpn

CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BIN_NAME} *.go
#upx -f --brute -o $BIN_NAME.upx $BIN_NAME
scp ${BIN_NAME} centos@${SUB_DOMAIN}.${CORP_DOMAIN}:

echo "${BIN_NAME} -enable-ses -email-bcc=support@${CORP_DOMAIN} -email-from=vpn@${CORP_DOMAIN}"
