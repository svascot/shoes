package shoes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamo *dynamodb.DynamoDB

type Shoe struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
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
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
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

}

func PutItem(shoe Shoe) {
	_, err := dynamo.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(shoe.Id)),
			},
			"Name": {
				S: aws.String(shoe.Name),
			},
			"Type": {
				S: aws.String(shoe.Type),
			},
			"Color": {
				S: aws.String(shoe.Color),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}
}

func CreateShoes() {
	var moreSneakers []Shoe

	var shoeOne Shoe
	shoeOne.Id = 1
	shoeOne.Color = "White"
	shoeOne.Name = "Jordan 1 low SE triple White"
	shoeOne.Type = "Jordan 1 low"

	moreSneakers = append(moreSneakers, shoeOne)

	var shoeTwo Shoe
	shoeTwo.Id = 2
	shoeTwo.Color = "Red"
	shoeTwo.Name = "Jordan 4 Metallic Red"
	shoeTwo.Type = "Jordan 4"

	moreSneakers = append(moreSneakers, shoeTwo)

	for _, s := range moreSneakers {
		PutItem(s)
	}

}

func GetShoes(w http.ResponseWriter, r *http.Request) {
	log.Println("GetShoes method")

	var moreSneakers []Shoe

	var shoeOne Shoe
	shoeOne.Color = "White"
	shoeOne.Name = "Jordan 1 low SE triple White"
	shoeOne.Type = "Jordan 1 low"

	moreSneakers = append(moreSneakers, shoeOne)

	var shoeTwo Shoe
	shoeTwo.Color = "Red"
	shoeTwo.Name = "Jordan 4 Metallic Red"
	shoeTwo.Type = "Jordan 4"

	moreSneakers = append(moreSneakers, shoeTwo)

	newJson, err := json.Marshal(moreSneakers)

	if err != nil {
		log.Println("Error marshalling json", err)
	}

	log.Println(newJson)
	w.Write(newJson)

}
