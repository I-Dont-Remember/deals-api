#! /bin/bash
# Generates a JS sdk of unknown quality
# https://docs.aws.amazon.com/apigateway/latest/developerguide/how-to-generate-sdk-cli.html

zip_file=./js-sdk.zip
# from the download
unzipped_dir=./apiGateway-js-sdk
sdk_dir=./js-sdk

if [ $# -ne 1 ]; then
    echo "Need Rest API id, for now it should be pvosoby5yl"
    exit 1
fi

if [ -d "$sdk_dir" ]; then 
    echo "SDK dir already exists: $sdk_dir, not downloading new sdk"
    exit 1
fi

aws apigateway get-sdk \
    --rest-api-id "$1" \
    --stage-name dev \
    --sdk-type javascript \
    "$zip_file" || exit 1

echo "Succesfully downloaded; unzip it and remove zip"
unzip "$zip_file"
rm "$zip_file"

mv "$unzipped_dir" "$sdk_dir"

echo "Moved sdk to dir: $sdk_dir"