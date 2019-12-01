package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

type PersonReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Person struct {
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

func getPersons(table dynamo.Table) (events.APIGatewayProxyResponse, error) {
	var persons []Person

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

func addPerson(table dynamo.Table, reqBody string) (events.APIGatewayProxyResponse, error) {

	id, _ := uuid.NewUUID()

	jsonBytes := []byte(reqBody)

	// Bodyを構造体に変換
	personReq := new(PersonReq)
	if err := json.Unmarshal(jsonBytes, personReq); err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("decode json error: %s", err.Error()))
	}

	// 書き込むための構造体を作成
	person := Person{
		Id:        id.String(),
		LastName:  personReq.LastName,
		FirstName: personReq.FirstName,
	}

	err := table.Put(person).Run()
	if err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("add person error: %s", err.Error()))
	}

	res, err := json.Marshal(person)
	if err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("create json error: %s", err.Error()))
	}

	return createRes(http.StatusCreated, string(res))
}

func delPerson(table dynamo.Table, id string) (events.APIGatewayProxyResponse, error) {
	err := table.Delete("Id", id).Run()
	if err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("delete person error: %s", err.Error()))
	}

	return createRes(http.StatusNoContent, "")
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// AWS SDKのセッション作成でエラーが発生した場合の処理
	if awsSess.Err != nil {
		return createRes(http.StatusInternalServerError,
			fmt.Sprintf("create aws session error: %s", awsSess.Err.Error()))
	}

	ddb := dynamo.New(awsSess.Sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	switch {
	// GET /persons
	case req.HTTPMethod == http.MethodGet && req.Path == "/persons":
		return getPersons(table)
	// POST /persons
	case req.HTTPMethod == http.MethodPost && req.Path == "/persons":
		return addPerson(table, req.Body)
	// DELETE /persons/{personId}
	case req.HTTPMethod == http.MethodDelete && req.Path ==
		fmt.Sprintf("/persons/%s", req.PathParameters["personId"]):
		return delPerson(table, req.PathParameters["personId"])
	}

	return createRes(http.StatusNotFound, "not found path or not allowed method")
}

func main() {
	lambda.Start(Handler)
}
