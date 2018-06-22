#! /bin/bash

GOOS=linux go build -o test test.go \
    && zip deployment.zip test \
    && aws lambda update-function-code --function-name Test --zip-file fileb://./deployment.zip

