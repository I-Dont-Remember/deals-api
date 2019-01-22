#! /bin/bash
command="awslocal"
dealsTable="Deals"
locationsTable="Locations"
campusesTable="Campuses"
analyticsTable="Analytics"

while getopts "rh" opt; do
    case $opt in
        r)
            command="aws"
            echo "Creating tables in AWS"
        ;;
        h)
            echo "Usage: $0 [ -r ] [ -h ], -r to use aws instead of awslocal"
            exit 0
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

# TODO: long-term todo, learn how to use Dynamo better with secondary indexes & whatnot
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

echo "Creating $campusesTable table..."
$command dynamodb create-table --table-name "$campusesTable" \
    --attribute-definitions AttributeName=slug,AttributeType=S \
    --key-schema AttributeName=slug,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

echo "Creating $analyticsTable table..."
$command dynamodb create-table --table-name "$analyticsTable" \
    --attribute-definitions AttributeName=timestamp,AttributeType=S \
    --key-schema AttributeName=timestamp,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

