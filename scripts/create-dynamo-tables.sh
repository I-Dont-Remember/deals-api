#! /bin/bash
command="awslocal"
dealsTable="Deals"
locationsTable="Locations"
campusesTable="Campuses"
analyticsTable="Analytics"


createDeals() {
    echo "Creating $1 table..."
    $command dynamodb create-table --table-name "$1" \
        --attribute-definitions AttributeName=id,AttributeType=S \
        --key-schema AttributeName=id,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
}

createLocations() {
    echo "Creating $1 table..."
    $command dynamodb create-table --table-name "$1" \
        --attribute-definitions AttributeName=id,AttributeType=S \
        --key-schema AttributeName=id,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
}

createCampuses() {
    echo "Creating $1 table..."
    $command dynamodb create-table --table-name "$1" \
        --attribute-definitions AttributeName=slug,AttributeType=S \
        --key-schema AttributeName=slug,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
}

createAnalytics() {
    echo "Creating $1 table..."
    $command dynamodb create-table --table-name "$1" \
        --attribute-definitions AttributeName=timestamp,AttributeType=S \
        --key-schema AttributeName=timestamp,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
}

while getopts "rhd" opt; do
    case $opt in
        r)
            command="aws"
            echo "Creating tables in AWS"
        ;;
        d)
            dev=true
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
createCampuses "$campusesTable"
createLocations "$locationsTable"
createDeals "$dealsTable"
createAnalytics "$analyticsTable"

if [ ! -z $dev ]; then
    createCampuses "$campusesTable"-dev
    createLocations "$locationsTable"-dev
    createDeals "$dealsTable"-dev
    createAnalytics "$analyticsTable"-dev
fi
