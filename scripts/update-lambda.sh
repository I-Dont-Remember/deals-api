#! /bin/bash
. ./helper.sh

command="awslocal"
# Assume no remote first, otherwise change it
lambdaName=$1
dirPath=$2
if [ $# -ne 0 ]; then
    if [ "$1" = "--remote" ]; then
        command="aws"
        lambdaName=$2
        dirPath=$3
    fi
else
    lambdaUsage
fi

# Check value isn't unset or empty string
if [ -z "$lambdaName" ]; then
    lambdaUsage
fi

if [ -z "$dirPath" ]; then
    lambdaUsage
fi

echo "Building and zipping function..."
cd "$dirPath" || echo "[!] couldn't cd to $dirPath." || exit 1
GOOS=linux go build -o main
zip deployment.zip main

echo "Updating function $lambdaName..."
$command lambda update-function-code --function-name "$lambdaName" \
    --zip-file fileb://deployment.zip

echo "Cleaning up repo..."
rm main deployment.zip
