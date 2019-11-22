package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	ddb := dynamo.New(sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	id := req.PathParameters["personId"]

	err = table.Delete("Id", id).Run()
	if err != nil {
		panic(err)
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: "",
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
