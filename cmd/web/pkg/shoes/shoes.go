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
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Color []string `json:"color"`
	Brand string   `json:"brand"`
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
			"Brand": {
				S: aws.String(shoe.Brand),
			},
			"Color": {
				SS: aws.StringSlice(shoe.Color),
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

func createMockShoes() {
	var moreSneakers []Shoe

	var shoeOne Shoe
	shoeOne.Id = 1
	shoeOne.Name = "Cement True Blue"
	shoeOne.Brand = "Nike"
	shoeOne.Color = []string{"Blue", "Gray", "White"}
	shoeOne.Type = "Air Jordan 1 Mid"

	var shoeTwo Shoe
	shoeTwo.Id = 2
	shoeTwo.Name = "Inside Out"
	shoeTwo.Brand = "Nike"
	shoeTwo.Color = []string{"Cream", "White", "Gray"}
	shoeTwo.Type = "Air Jordan 1 Low SE Craft"

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
		PutItem(s)
	}

}

func CreateMockData(w http.ResponseWriter, r *http.Request) {
	createMockShoes()
	w.Write([]byte("Mock data created"))

}

func GetShoes(w http.ResponseWriter, r *http.Request) {

	tbl, err := dynamo.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	newJson, err := json.Marshal(tbl.Items)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(newJson)

}
