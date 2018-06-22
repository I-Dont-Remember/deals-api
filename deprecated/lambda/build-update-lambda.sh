#! /bin/bash

GOOS=linux go build -o main main.go \
    && zip deployment.zip main \
    && aws lambda update-function-code --function-name Deals --zip-file fileb://./deployment.zip

