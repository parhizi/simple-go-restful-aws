package main

import (
	"testing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"encoding/json"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws"
)

type TestCase struct {
	Name                        string
	Request             		events.APIGatewayProxyRequest
	//InputId                     events.APIGatewayProxyRequest
	//InputIdString               string
	inputedItems				map [string]* dynamodb.AttributeValue
	//DatabaseOutput              dynamodb.PutItemOutput
	ExpectedBody                string
	ExpectedStatusCode          int
	ExpectedDatabaseOutput      dynamodb.PutItemOutput
	ExpectedError               error
}

// Mocking DynamoDB through dynamodbiface.
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	// Other return values expected to store, i.e: "payload map[string]string" or "err error"
}

// Custom PutItem function for overriding the PutItem of getDeviceById.go for using in test scenarios.
// Mocking PutItem output to the a desire valid response.
func (self *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	MockOutput := new(dynamodb.PutItemOutput)
	return MockOutput, nil
}

// Put function in addDevice.go signature: input: (item map[string] *dynamodb.AttributeValue), output: (*dynamodb.PutItemOutput, error)
func TestPut(t *testing.T) {
	// Preparing a DynamoDB PuItemOutput data type as expected from a DB response.
	MockInput := dynamodb.PutItemInput{}
	MockInput.SetItem(
		map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String("id1")},
			"deviceModel": &dynamodb.AttributeValue{S: aws.String("testDeviceModel")},
			"name": &dynamodb.AttributeValue{S: aws.String("testName")},
			"note": &dynamodb.AttributeValue{S: aws.String("testNote")},
			"serial": &dynamodb.AttributeValue{S: aws.String("testSerial")},
		},
	)

	//When  we have come up to putting item on DB, the Body data is standard and without any issues,
	//Because we have validated the input body request in ValidateInputs functions beforehand in addDevice.go file
	testCase := TestCase{
		Name:            "** Testing JSON with proper fields **",
		inputedItems:    MockInput.Item,
		ExpectedError:	 nil,
	}

	// Prepare AWS & DynamoDB session for mocking.
	test_aws := new(AmazonWebServices)
	test_aws.DynamoDB = &MockDynamoDB{}

	_, err := test_aws.Put(testCase.inputedItems)

	//Function here is %100 proof, so no error will happen.
	if err != testCase.ExpectedError {
		t.Errorf("%s \n \t<expected error: %t> <resulted error: %t>", testCase.ExpectedError, err)
	}
} // End of TestPut function.

// ValidateDatabaseResult function in addDevice.go signature: input: (request events.APIGatewayProxyRequest), output: (Device, error)
func TestAddDevice(t *testing.T) {
	testCases := []TestCase{
		{
			Name:                "** Testing: Empty body input. **",
			Request:             events.APIGatewayProxyRequest{Body: ""},
			ExpectedBody:         "No inputs provided, please provide inputs in JSON format.",
			ExpectedStatusCode:  400,
		},
		{
			Name:                "** Testing: Wrong JSON format. **",
			Request:             events.APIGatewayProxyRequest{Body: "{{{}"},
			ExpectedBody:         "Wrong format: Inputs must be a valid JSON.",
			ExpectedStatusCode:  400,
		},
		{
			Name:                "** Testing: JSON with missing field - ID **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"testName\" , \"note\":\"testNote\" , \"serial\":\"testSerial\" }"},
			ExpectedBody:         "Missing field: ID",
			ExpectedStatusCode:  400,
		},
		{
			Name:                "** Testing: JSON with missing field - Device Model **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"\" , \"name\":\"testName\" , \"note\":\"testNote\" , \"serial\":\"testSerial\" }"},
			ExpectedBody:         "Missing field: Device Model",
			ExpectedStatusCode:  400,
		},

		{
			Name:                "** Testing: JSON with missing field - Name **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"\" , \"note\":\"testNote\" , \"serial\":\"testSerial\" }"},
			ExpectedBody:         "Missing field: Name",
			ExpectedStatusCode:  400,
		},

		{
			Name:                "** Testing: JSON with missing field - Note **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"testName\" , \"note\":\"\" , \"serial\":\"testSerial\" }"},
			ExpectedBody:         "Missing field: Note",
			ExpectedStatusCode:  400,
		},

		{
			Name:                "** Testing: JSON with missing field - Serial **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"testName\" , \"note\":\"testNote\" , \"serial\":\"\" }"},
			ExpectedBody:         "Missing field: Serial",
			ExpectedStatusCode:  400,
		},
		{	// In Testing environment, as we don't access AWS's OS environment variable and other real world parameters, can not reach to
			// HTTP code 201 point in here, unless we prepare a mock server for it.
			Name:                "** Testing: JSON with proper fields. **",
			Request:             events.APIGatewayProxyRequest{Body: "{\"id\":\"1\",\"deviceModel\":\"testDeviceModel\",\"name\":\"testName\",\"note\":\"testNote\",\"serial\":\"testSerial\"}"},
			ExpectedBody:        "Internal Server Error\nDatabase error.",
			//ExpectedBody:        "{\"id\":\"1\",\"deviceModel\":\"testDeviceModel\",\"name\":\"testName\",\"note\":\"testNote\",\"serial\":\"testSerial\"}" ,
			ExpectedStatusCode:  500, //201
		},
	}


	for _, test := range testCases {
		// Executing each test cases scenario.
		response, _ := AddDevice(test.Request)
		if response.StatusCode != test.ExpectedStatusCode || response.Body != test.ExpectedBody{
			t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
		}
	}

} // end of TestAddDevice function
