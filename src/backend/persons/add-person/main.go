package main

import (
	"encoding/json"
	"io"
	"log"
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

type Params struct {
	sess      *session.Session
	req       events.APIGatewayProxyRequest
	outStream io.Writer
	errStream io.Writer
}

func Run(params *Params) (events.APIGatewayProxyResponse, error) {

	log.SetOutput(params.errStream)

	//svc := dynamodb.New(params.sess)

	reqBody := params.req.Body

	jsonBytes := ([]byte)(reqBody)
	personReq := new(PersonReq)
	if err := json.Unmarshal(jsonBytes, personReq); err != nil {
		log.Fatal(err)
	}

	ddb := dynamo.New(params.sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	table.Put(Person{
		Id:        id.String(),
		FirstName: personReq.FirstName,
		LastName:  personReq.LastName,
	})
	//
	//putParams := &dynamodb.PutItemInput{
	//	TableName: aws.String(os.Getenv("TABLE_NAME")),
	//	Item: map[string]*dynamodb.AttributeValue{
	//		"Id": {
	//			S: aws.String(id.String()),
	//		},
	//		"FirstName": {
	//			S: aws.String(personReq.FirstName),
	//		},
	//		"LastName": {
	//			S: aws.String(personReq.LastName),
	//		},
	//	},
	//}

	//putItem, err := svc.PutItem(putParams)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(putItem)

	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: "",
	}

	return resp, nil
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	params := &Params{
		sess:      sess,
		req:       req,
		outStream: os.Stdin,
		errStream: os.Stdout,
	}

	resp, err := Run(params)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
