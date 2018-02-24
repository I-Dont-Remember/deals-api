package awshelper

import (
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
     //backslash might be necessary 
     //https://stackoverflow.com/questions/33801460/dynamo-db-local-connection-refused#37106007
     localUrl = "http://localhost:8000/"
    region = "us-east-2"
    Tables = map[string]*dynamodb.CreateTableInput {
        "Deals": {
            TableName: aws.String("Deals"),
            AttributeDefinitions: []*dynamodb.AttributeDefinition{
                {
                    AttributeName: aws.String("Id"),
                    AttributeType: aws.String("S"),
                },
            },
            KeySchema: []*dynamodb.KeySchemaElement{
                {
                    AttributeName: aws.String("Id"),
                    KeyType: aws.String("HASH"),
                },
            },
            ProvisionedThroughput: &dynamodb.ProvisionedThroughput {
                ReadCapacityUnits: aws.Int64(1),
                WriteCapacityUnits: aws.Int64(1),
            },
        },
        "SuggestedDeals": {
            TableName: aws.String("SuggestedDeals"),
        },
    }
)

// getSVC returns the svc object used to access the DynamoDB service
func getSVC(local bool) *dynamodb.DynamoDB {
    sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if local {
        return dynamodb.New(sess, aws.NewConfig().WithEndpoint(localUrl))
    }
    return dynamodb.New(sess)
}

// CreateTable creates the DynamoDB table in either local or AWS
func CreateTable(local bool, name string) error {
    if Tables[name] == nil {
        return fmt.Errorf("[x] awshelper: no table entry available for %s", name)
    }

    fmt.Printf("[*] trying to create %s...\n", name)
    input := Tables[name]
    svc := getSVC(local)
    _, err := svc.CreateTable(input)
    if err != nil {
        fmt.Println("[x] failed creating table...")
        fmt.Println(err.Error())
    } else {
        fmt.Printf("[+] created table %s...\n", name)
    }

    return nil
}
