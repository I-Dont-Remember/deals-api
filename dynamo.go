package main

import (
    "os"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
    local_env = "LOCAL_DYNAMODB"
    region = "us-east-2"
)


func check(e error) {
    if e != nil {
        fmt.Printf("Error: %s\n", e)
        os.Exit(2)
    }
}


func useAWS(sess *session.Session) *dynamodb.DynamoDB {
    fmt.Println("Using AWS DynamoDB connection...")
    return dynamodb.New(sess)
}


func useLocal(sess *session.Session) *dynamodb.DynamoDB {
    fmt.Println("Using local DynamoDB connection...")
    //backslash might be necessary https://stackoverflow.com/questions/33801460/dynamo-db-local-connection-refused#37106007
    return dynamodb.New(sess, aws.NewConfig().WithEndpoint("http://localhost:8000/"))
}


func main() {
    var svc *dynamodb.DynamoDB
    sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
    check(err)

    if os.Getenv(local_env) == "1" {
        svc = useLocal(sess)
    } else {
        svc = useAWS(sess)
    }

    input := &dynamodb.ListTablesInput{}

    result, err := svc.ListTables(input)
    if err != nil {
        fmt.Println(err.Error())
    }

    fmt.Println(result)

    t_input := &dynamodb.DescribeTableInput{
                    TableName: aws.String(*result.TableNames[0]),
            }

    t_result, err := svc.DescribeTable(t_input)
    if err != nil {
        fmt.Println(err.Error())
    }

    fmt.Println(t_result)
}
