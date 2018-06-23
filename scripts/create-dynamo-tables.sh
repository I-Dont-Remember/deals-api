#! /bin/bash
command="awslocal"
dealsTable="devDeals"
locationsTable="devLocations"

while getopts ":d:l:r" opt; do
    case $opt in
        d)
            dealsTable=$OPTARG
            echo "-d deals table $dealsTable"
        ;;
        l)
            locationsTable=$OPTARG
            echo "-l Location table $locationsTable"
        ;;
        r)
            command="aws"
            echo "-r remote $command"
        ;;
        \?)
            echo "Invalid Option: -$OPTARG"
            exit 1
        ;;
        :)
            echo "Option -$OPTARG requires an argument."
            exit 1
    esac
done

# Understanding what to do for tables https://www.dynamodbguide.com/working-with-multiple-items
# "Give me all of the deals(Range) from a particular location(Hash)"
# This then doesn't work,because primary id is now location, which has many so you only get 1 deal per location
echo "Creating $dealsTable table..."
$command dynamodb create-table --table-name "$dealsTable" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

echo "Creating $locationsTable table..."
$command dynamodb create-table --table-name "$locationsTable" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
