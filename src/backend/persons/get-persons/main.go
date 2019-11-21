package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
)

type PersonRes struct {
	Id        string `dynamo:"Id"`
	FirstName string `dynamo:"FirstName"`
	LastName  string `dynamo:"LastName"`
}

func Handler() (events.APIGatewayProxyResponse, error) {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	ddb := dynamo.New(sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	var persons []PersonRes

	err = table.Scan().All(&persons)
	if err != nil {
		panic(err)
	}
	fmt.Println(persons)

	jsonBytes, _ := json.Marshal(persons)

	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(jsonBytes),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
