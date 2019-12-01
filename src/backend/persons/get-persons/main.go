package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type PersonRes struct {
	Id        string `dynamo:"Id"`
	FirstName string `dynamo:"FirstName"`
	LastName  string `dynamo:"LastName"`
}

type AwsSess struct {
	Sess *session.Session `json:"sesson"`
	Err  error            `json:"error"`
}

var awsSess AwsSess

func init() {
	awsSess.Sess, awsSess.Err = session.NewSession()
	if awsSess.Err != nil {
		panic(awsSess.Err)
	}
}

func createRes(code int, msg string) (events.APIGatewayProxyResponse, error) {
	header := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}
	res := events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    header,
		Body:       msg,
	}
	return res, nil
}

func Handler() (events.APIGatewayProxyResponse, error) {

	// AWS SDKのセッション作成でエラーが発生した場合の処理
	if awsSess.Err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("create aws session error: %s", awsSess.Err.Error()))
	}

	ddb := dynamo.New(awsSess.Sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	var persons []PersonRes

	err := table.Scan().All(&persons)
	if err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("scan error: %s", err.Error()))
	}

	jsonBytes, err := json.Marshal(persons)
	if err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("create json error: %s", err.Error()))
	}

	return createRes(http.StatusOK, string(jsonBytes))
}

func main() {
	lambda.Start(Handler)
}
