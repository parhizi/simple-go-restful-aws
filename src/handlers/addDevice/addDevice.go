package main

import (
	"os"
	"errors"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"fmt"
)

type Device struct {
	ID          string  `json:"id"`
	DeviceModel string  `json:"deviceModel"`
	Name        string  `json:"name"`
	Note  		string  `json:"note"`
	Serial   	string  `json:"serial"`
}

type AmazonWebServices struct {
	Config *aws.Config
	Session *session.Session
	DynamoDB dynamodbiface.DynamoDBAPI
}

// Preparing AWS & DynamoDB session
func ConfigureAws()(*AmazonWebServices) {
	region := os.Getenv("AWS_REGION")
	var Aws *AmazonWebServices = new(AmazonWebServices)
	Aws.Config = &aws.Config{Region: aws.String(region),}
	var err error
	Aws.Session, err = session.NewSession(Aws.Config)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		var svc *dynamodb.DynamoDB = dynamodb.New(Aws.Session)
		Aws.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	}
	return Aws
}

// Preparing DynamoDB Session and Calling DB's PutItem function inside.
func (self *AmazonWebServices) Put(item map[string] *dynamodb.AttributeValue) (*dynamodb.PutItemOutput, error) {
	// Get table name from OS's environment
	tableName := aws.String(os.Getenv("DEVICES_TABLE_NAME"))
	var input = &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}
	// Calling either PutItem function of interface, defined in addDevice_test.go file, or api with the input we've provided.
	// In mock case, the PutItem function of getDeviceById_test.go will be called(interface.go)
	// In real deployment environment, the PutItem function of aws (api.go) will be called.
	result, err := self.DynamoDB.PutItem(input)
	// Todo: ignoring the err! Not Funny!
	return result, err
}

// The handler function which will be first started from main function.
func AddDevice(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// First & foremost we have to validate user input.
	NewDevice, err := ValidateInputs(request)
	// if inputs are not suitable, return HTTP error code 400.
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "" + err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Prepare a new AWS & DynamoDB session and configure it.
	TestAws := ConfigureAws()

	// // Serialization/Encoding "NewDevice" in "item" for using in DynamoDB functions.
	item, _ := dynamodbattribute.MarshalMap(NewDevice)

	// Till now the user have provided a valid data input.
	// Let's add it to the DynamoDB talble.
	_, err = TestAws.Put(item)

	// If internal database errors occurred, return HTTP error code 500.
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error\nDatabase error.",
			StatusCode: 500,
		}, nil
	}

	// Serialization/Encoding "NewDevice" to JSON.
	jsonResponse, _ := json.Marshal(NewDevice)
	return events.APIGatewayProxyResponse {
		Body: string(jsonResponse),
		// Everthing looks fine, return HTTP 201
		StatusCode: 201,
	}, nil
} // End of AddDevice function

func ValidateInputs(request events.APIGatewayProxyRequest) (Device, error) {

	NewDevice := Device {}
	ErrorMessage := ""

	if len(request.Body) == 0 {
		ErrorMessage = "No inputs provided, please provide inputs in JSON format."
		return Device{}, errors.New(ErrorMessage)
	}

	// De-serialize "request.Body" which is in JSON format into "NewDevice" in Go object.
	var err = json.Unmarshal([]byte(request.Body), &NewDevice)

	if err != nil {
		ErrorMessage = "Wrong format: Inputs must be a valid JSON."
		return Device{}, errors.New(ErrorMessage)
	}

	if len(NewDevice.ID) == 0 {
		ErrorMessage = "Missing field: ID"
		return Device{}, errors.New(ErrorMessage)
	}

	if len(NewDevice.DeviceModel) == 0 {
		ErrorMessage = "Missing field: Device Model"
		return Device{}, errors.New(ErrorMessage)
	}

	if len(NewDevice.Name) == 0 {
		ErrorMessage = "Missing field: Name"
		return Device{}, errors.New(ErrorMessage)
	}
	if len(NewDevice.Note) == 0 {
		ErrorMessage = "Missing field: Note"
		return Device{}, errors.New(ErrorMessage)
	}

	if len(NewDevice.Serial) == 0 {
		ErrorMessage = "Missing field: Serial"
		return Device{}, errors.New(ErrorMessage)
	}

	// Everything looks fine, return created NewDevice in Go struct.
	return NewDevice, nil
} // End of ValidateInputs function.

func main() {
	lambda.Start(AddDevice)
}
