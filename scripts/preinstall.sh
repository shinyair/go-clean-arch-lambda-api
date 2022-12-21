#!/bin/sh

# Serverless Framework
echo ">> install serverless framework"
npm install -g serverless
serverless --version
echo "<< done"

# golangci-lint
echo ">> install golangci-lint"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.48.0
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0
golangci-lint --version
echo "<< done"

# swaggo
echo ">> install swaggo"
go install github.com/swaggo/swag/cmd/swag@latest
which swag
echo "<< done"