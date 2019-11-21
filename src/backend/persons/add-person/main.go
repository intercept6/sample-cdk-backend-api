package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type PersonReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	id, _ := uuid.NewUUID()

	reqBody := event.Body
	fmt.Println(reqBody)

	jsonBytes := ([]byte)(reqBody)
	personReq := new(PersonReq)
	if err := json.Unmarshal(jsonBytes, personReq); err != nil {
		panic(err)
	}

	fmt.Println("Id: ", id)
	fmt.Println("FirstName: ", personReq.FirstName)
	fmt.Println("LastName: ", personReq.LastName)

	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id.String()),
			},
			"FirstName": {
				S: aws.String(personReq.FirstName),
			},
			"LastName": {
				S: aws.String(personReq.LastName),
			},
		},
	}

	putItem, err := svc.PutItem(putParams)
	if err != nil {
		panic(err)
	}
	fmt.Println(putItem)

	resp := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body:            "",
		IsBase64Encoded: false,
	}

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
