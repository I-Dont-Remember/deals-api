#! /bin/bash
command="awslocal"
if [ $# -ne 0 ]; then
    if [ $1 = "--remote" ]; then
        command="aws"
    fi
fi

echo "Creating Deals table..."
$command dynamodb create-table --table-name Deals \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

echo "Creating Locations table..."
$command dynamodb create-table --table-name Locations \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
