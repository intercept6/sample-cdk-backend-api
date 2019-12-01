package main

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	dynamodbEndpoint = "http://localhost:4569"
	region           = "ap-northeast-1"
)

func createSess(t *testing.T) AwsSess {
	t.Helper()
	awsSess = AwsSess{}
	awsSess.Sess, awsSess.Err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("DUMMY", "DUMMY", "DUMMY"),
		Endpoint:    aws.String(dynamodbEndpoint),
		Region:      aws.String(region),
	})
	if awsSess.Err != nil {
		t.Fatal(awsSess.Err)
	}
	return awsSess
}

type Person struct {
	Id        string `json:"Id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

func createTable(t *testing.T, ddb *dynamo.DB) string {
	t.Helper()
	type PersonsTable struct {
		Id string `dynamo:"Id,hash"`
	}
	res := ddb.CreateTable("test_table_for_get_persons", PersonsTable{})
	if err := res.Run(); err != nil {
		t.Fatal(err)
	}
	return "test_table"
}

func createRecord(t *testing.T, ddb *dynamo.DB, tableName string) []Person {
	t.Helper()
	table := ddb.Table(tableName)
	item := Person{Id: "556350d2-e993-4fb9-8242-c496a0664bb3", LastName: "Taro", FirstName: "Yamada"}
	err := table.Put(item).Run()
	if err != nil {
		t.Fatal(err)
	}

	return []Person{item}
}

func deleteTable(t *testing.T, ddb *dynamo.DB, tableName string) {
	t.Helper()

	table := ddb.Table(tableName)
	table.DeleteTable()
}

func TestGetPerson(t *testing.T) {
	awsSess = createSess(t)
	ddb := dynamo.New(awsSess.Sess)
	tableName := createTable(t, ddb)
	defer deleteTable(t, ddb, tableName)
	wantBody := createRecord(t, ddb, tableName)
	if err := os.Setenv("TABLE_NAME", tableName); err != nil {
		t.Fatal(err)
	}

	// Handlerは必ずerrを返さない
	res, _ := Handler()

	var gotBody []Person
	if err := json.Unmarshal([]byte(res.Body), &gotBody); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got: %v, want: %v", res.StatusCode, http.StatusOK)
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Errorf("got: %v, want: %v", res.Body, wantBody)
	}
}
