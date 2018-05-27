#! /bin/bash
command="awslocal"
# Assume no remote first, otherwise change it
lambdaName=$1
zipPath=$2
if [ $# -ne 0 ]; then
    if [ $1 = "--remote" ]; then
        command="aws"
        lambdaName=$2
        zipPath=$3
    fi
else
    echo "Usage: script [--remote] lambdaName zipFile"
    exit 1
fi

# Check value isn't unset or empty string
if [ -z "$lambdaName" ]; then
    echo "Usage: script [--remote] lambdaName zipFile"
    exit 1
fi

if [ -z "$zipPath" ]; then
    echo "Usage: script [--remote] lambdaName zipFile"
    exit 1
fi

echo "Creating function $lambdaName..."
$command lambda create-function --function-name $lambdaName \
    --runtime go1.x \
    --role DealsLambda \
    --handler main \
    --zip-file fileb://$zipPath