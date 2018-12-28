#! /bin/bash

# https://stackoverflow.com/questions/33667334/exporting-api-definition-from-aws-api-gateway

if [ $# -ne 1 ]; then
    echo "Need Rest API id"
    exit 1
fi

aws apigateway get-export \
    --rest-api-id "$id" \
    --stage-name dev \
    --export-type swagger \
    ./swagger.json