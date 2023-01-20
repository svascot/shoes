package shoes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

var dynamo *dynamodb.DynamoDB

type Shoe struct {
	Id    int         `json:"id"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Color []string    `json:"color"`
	Brand string      `json:"brand"`
	Dates []time.Time `json:"dates"`
}

const TABLE_NAME = "qa_shoes"

func init() {
	dynamo = connectDynamo()
}

func connectDynamo() (db *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))
}

func CreateTable() {

	log.Println("Creating shoes table.")

	// table schema
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(TABLE_NAME),
	}

	//creates the table
	_, err := dynamo.CreateTable(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
	log.Println("Shoes table created.")

}

func GetShoe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shoe := getShoe(params["id"])
	json.NewEncoder(w).Encode(shoe)

}

func DeleteShoe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
	json.NewEncoder(w).Encode(id)
}

func AddShoe(w http.ResponseWriter, r *http.Request) {
	var shoe Shoe
	json.NewDecoder(r.Body).Decode(&shoe)
	putShoe(shoe)
	json.NewEncoder(w).Encode(shoe)
}

func WearToday(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shoe := getShoe(params["id"])
	shoe.Dates = append(shoe.Dates, time.Now())
	putShoe(shoe)
	json.NewEncoder(w).Encode(shoe)
}

func putShoe(shoe Shoe) {

	item, err := dynamodbattribute.MarshalMap(shoe)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func getShoe(id string) Shoe {
	var shoe Shoe

	result, err := dynamo.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	dynamodbattribute.UnmarshalMap(result.Item, &shoe)

	return shoe
}

func GetAllShoes(w http.ResponseWriter, r *http.Request) {

	tbl, err := dynamo.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	var shoes []Shoe

	for _, dShoe := range tbl.Items {
		var shoe Shoe
		err = dynamodbattribute.UnmarshalMap(dShoe, &shoe)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		shoes = append(shoes, shoe)
	}

	newJson, err := json.Marshal(shoes)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(newJson)
}

// Mock Data
func createMockShoes() {
	fmt.Println("creating mock shoes")
	var moreSneakers []Shoe

	var shoeOne Shoe
	shoeOne.Id = 1
	shoeOne.Name = "Cement True Blue"
	shoeOne.Brand = "Nike"
	shoeOne.Color = []string{"Blue", "Gray", "White"}
	shoeOne.Type = "Air Jordan 1 Mid"
	shoeOne.Dates = append(shoeOne.Dates, time.Now())

	var shoeTwo Shoe
	shoeTwo.Id = 2
	shoeTwo.Name = "Inside Out"
	shoeTwo.Brand = "Nike"
	shoeTwo.Color = []string{"Cream", "White", "Gray"}
	shoeTwo.Type = "Air Jordan 1 Low SE Craft"
	shoeTwo.Dates = []time.Time{time.Now()}

	var shoeThree Shoe
	shoeThree.Id = 3
	shoeThree.Name = "Adidas Ultraboost"
	shoeThree.Brand = "Adidas"
	shoeThree.Color = []string{"Black"}
	shoeThree.Type = "Ultraboost"

	var shoeFour Shoe
	shoeFour.Id = 4
	shoeFour.Name = "Adidas Grand Court"
	shoeFour.Brand = "Adidas"
	shoeFour.Color = []string{"Black", "White"}
	shoeFour.Type = "Grand Court"

	moreSneakers = append(moreSneakers, shoeOne)
	moreSneakers = append(moreSneakers, shoeTwo)
	moreSneakers = append(moreSneakers, shoeThree)
	moreSneakers = append(moreSneakers, shoeFour)

	for _, s := range moreSneakers {
		putShoe(s)
		fmt.Sprintln("shoe added")
	}
}

func CreateMockData(w http.ResponseWriter, r *http.Request) {
	createMockShoes()
	w.Write([]byte("Mock data created"))
}
