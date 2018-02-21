#! /bin/bash

function makeTable {
    aws dynamodb create-table \
        --table-name Deals \
        --attribute-definitions \
            AttributeName=Id,AttributeType=S \
        --key-schema AttributeName=Id,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
}


if $(hash aws &> /dev/null); then
    echo "Making new DynamoDB table..."
    makeTable
else
    echo "'aws' is not an existing command on the system..."
fi

#        --attribute-definitions \
#            AttributeName=Id,AttributeType=S \
#            AttributeName=Day,AttributeType=S \
#            AttributeName=Location,AttributeType=S \
#            AttributeName=Deal,AttributeType=S \
#        --key-schema AttributeName=Id,KeyType=HASH AttributeName=Day,KeyType=RANGE AttributeName=Location,KeyType=RANGE AttributeName=Deal,KeyType=RANGE \

