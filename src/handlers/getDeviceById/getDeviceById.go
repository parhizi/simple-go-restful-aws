package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"os"
	"types"
)

type AmazonWebServices struct {
	Config   *aws.Config
	Session  *session.Session
	DynamoDB dynamodbiface.DynamoDBAPI
}

// Prepare a new AWS & DynamoDB session, then configure it.
var TestAws *AmazonWebServices

func init() {
	region := os.Getenv("AWS_REGION")
	var Aws *AmazonWebServices = new(AmazonWebServices)
	Aws.Config = &aws.Config{Region: aws.String(region)}
	var err error
	Aws.Session, err = session.NewSession(Aws.Config)
	if err != nil {
		// Logs error on Amazon CloudWatch. It's sysadmin's duty to handle it.
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		var svc *dynamodb.DynamoDB = dynamodb.New(Aws.Session)
		Aws.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	}
	// Instantiate a global session in TestAws
	TestAws = Aws
}

// Preparing DynamoDB Session and Calling DB's GetItem function inside.
func (self *AmazonWebServices) Get(id string) (*dynamodb.GetItemOutput, error) {
	// Get desire table's name from OS's environmental varible.
	tableName := aws.String(os.Getenv("DEVICES_TABLE_NAME"))

	// Putting tableName and the id which we have received previously from client side by GET method.
	var input = &dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	// Calling either GetItem function of interface, defined in getDeviceById_test.go file, or api with the input we've provided.
	// In mock case, the GetItem function of getDeviceById_test.go will be called(interface.api)
	// In real deployment environment, the GetItem function of aws (api.go) will be called.
	result, err := self.DynamoDB.GetItem(input)
	return result, err
}

// The handler function which will be first started from main function.
func GetDeviceById(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// The id which user has sent through GET method.
	id := request.PathParameters["id"]

	// If no id have been provided, return HTTP error code 404.
	if id == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Missing field : id",
			StatusCode: 404,
		}, nil
	}

	// Till now the user have provided an id in string type.
	// Let's see whether it's existed on DB or not.
	result, err := TestAws.Get(id)

	// Checking the result of the DynamoDB query.
	ValidationResult := ValidateDatabaseResult(result, err)

	// Return the result in ...
	return ValidationResult, nil
} // End of GetDeviceById function

func ValidateDatabaseResult(result *dynamodb.GetItemOutput, err error) events.APIGatewayProxyResponse {

	// If an internal error have occurred in the database, return HTTP error code 500.
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       string("Internal Server Error."),
			StatusCode: 500,
		}
	}

	// If no item have been founded, return HTTP error code 404.
	if len(result.Item) == 0 {
		return events.APIGatewayProxyResponse{
			Body:       string("Desired device not found."),
			StatusCode: 404,
		}
	}

	// Till now the input id have been founded.
	// Let's convert this founded "result.item" from DB which is in DynamoDB type to Go struct.
	item := types.Device{}
	// Deserialization/Decoding "result.Item" to Go struct.
	dynamodbattribute.UnmarshalMap(result.Item, &item)

	// Serialization/Encoding item to JSON.
	FoundedDeviceJson, _ := json.Marshal(item)

	// Return founded item as JSON type with 200 HTTP status code.
	return events.APIGatewayProxyResponse{
		Body:       string(FoundedDeviceJson),
		StatusCode: 200,
	}
} // End of ValidateDatabaseResult function

func main() {
	lambda.Start(GetDeviceById)
}
