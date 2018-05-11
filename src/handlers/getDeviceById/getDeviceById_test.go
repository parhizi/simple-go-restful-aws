package main

import (
    "errors"
    "testing"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type TestCase struct {
    Name                    string
    Request                 events.APIGatewayProxyRequest
    Input                   string
    MockDatabaseOutput      dynamodb.GetItemOutput
    ExpectedDatabaseOutput  dynamodb.GetItemOutput
    ExpectedBody            string
    ExpectedStatusCode      int
    Error                   error
}

// Mocking DynamoDB through dynamodbiface.
type MockDynamoDB struct {
    dynamodbiface.DynamoDBAPI
    // Other return values expected to store, i.e: "payload map[string]string" or "err error"
}

// Custom GetItem function for overriding the GetItem of getDeviceById.go for using in test scenarios.
// Mocking GetItem output to the a desire valid response.
func (self *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (output *dynamodb.GetItemOutput,err error) {
    mockOutput := new(dynamodb.GetItemOutput)
    inputID := input.Key["id"].S

    // Checking whether the test case id input is equal to the mocked DB's id value or not.
    if *inputID == "id_test" {
        mockOutput.SetItem(
            // Setting mocked values.
            map[string]*dynamodb.AttributeValue{
                "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
                "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
                "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
                "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
                "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
            },
        )
    }
    return mockOutput, err
}


// Get function in getDeviceById.go signature: input: (id string) , output: (*dynamodb.GetItemOutput, error)
func TestGet(t *testing.T)  {

    // Preparing a DynamoDB GetItemOutput data type as expected from a DB response.
    MockOutput := dynamodb.GetItemOutput{}
    EmptyOutput := dynamodb.GetItemOutput{}
    // Setting mock items values
    MockOutput.SetItem(
        map[string]*dynamodb.AttributeValue{
            "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
            "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
            "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
            "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
            "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
        },
    )

    TestCases := []TestCase {
        {
            Name:                       "** Requested ID exists. **",
            Input:                      "id_test",
            // MockOutput is a well defined & set of DynamoDB values, as a successful DB return founded item is.
            ExpectedDatabaseOutput:     MockOutput,
        },
        {
            Name:                       "** Requested ID does not exist! **",
            Input:                      "NotExistedTestID",
            // EmptyOutput is just like a null DynamoDB return value for not founded item.
            ExpectedDatabaseOutput:     EmptyOutput,
        },
    }

    // Prepare AWS & DynamoDB session for mocking.
    test_aws := new(AmazonWebServices)
    test_aws.DynamoDB = &MockDynamoDB{}

    for _, test := range TestCases {
        // Executing each test cases scenario.
        response, _ := test_aws.Get(test.Input)
        if len(response.GoString()) != len(test.ExpectedDatabaseOutput.GoString()) {
            t.Errorf("%s \n \t<expected output: \n%s> \n<resulted output: \n%s>", test.Name, test.ExpectedDatabaseOutput.GoString(), response.GoString())
        }
    }
} // End of TestGet function

// GetDeviceById function in getDeviceById.go signature: input: (request events.APIGatewayProxyRequest), output: (events.APIGatewayProxyResponse, error)
func TestGetDeviceById(t *testing.T) {

    TestCases := []TestCase{
        {
            Name:                   "** Testing: Empty id input. **",
            Request:                events.APIGatewayProxyRequest{PathParameters: map[string]string {"id":""}},
            ExpectedBody:           "No input provided, please provide it in JSON format.",
            ExpectedStatusCode:     404,
        },

        {
            Name:                   "** Testing: Desire id does not exist. **",
            Request:                events.APIGatewayProxyRequest{PathParameters: map[string]string {"id":"doesn't existed"}},
            ExpectedBody:           "Wrong format: Input must be a valid JSON.",
            ExpectedStatusCode:     500,
        },

        //{
        //// In Testing environment, as we don't access AWS's OS environment variable and other real world parameters, can not reach to
        //// HTTP code 201 point in here, unless we prepare a mock server for it.
        //  Name:                   "** Testing: Proper id which does exist on DB. **",
        // 	Request:                events.APIGatewayProxyRequest{PathParameters: map[string]string {"id":"id1"}},
        //  //ExpectedBody:         "{\"id\":\"1\",\"deviceModel\":\"testDeviceModel\",\"name\":\"testName\",\"note\":\"testNote\",\"serial\":\"testSerial\"}" ,
        // 	ExpectedStatusCode:     201,
        //},

    }

    for _, test := range TestCases {
        // Executing each test cases scenario.
        response, _ := GetDeviceById(test.Request)

        if response.StatusCode != test.ExpectedStatusCode {
            t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d>", test.Name, test.ExpectedStatusCode, response.StatusCode)
        }
    }

} // End of TestGetDeviceById function

// ValidateDatabaseResult function in getDeviceById.go signature: input: (result *dynamodb.GetItemOutput, err error), output: (events.APIGatewayProxyResponse)
func TestValidateDatabaseResult(t *testing.T) {
    // Preparing a DynamoDB GetItemOutput data type as expected DB response.
    MockOutput := dynamodb.GetItemOutput{}
    EmptyOutput := dynamodb.GetItemOutput{}
    // Setting mock items values
    MockOutput.SetItem(
        map[string]*dynamodb.AttributeValue{
            "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
            "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
            "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
            "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
            "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
        },
    )

    TestCases := []TestCase{

        {
            Name:                   "** Database Unexpected Error **",
            Request:                events.APIGatewayProxyRequest{},
            MockDatabaseOutput:     EmptyOutput,
            Error:                  errors.New("unexpected Error has occurred"),
            ExpectedBody:           "Internal Server Error.",
            ExpectedStatusCode:     500,
        },

        {
            Name:                   "** Database Returns Empty Result **",
            Request:                events.APIGatewayProxyRequest{},
            MockDatabaseOutput:     EmptyOutput,
            ExpectedBody:           "Desired device not found.",
            ExpectedStatusCode:     404,
        },

        {
            Name:                   "** Database Returns founded device **",
            MockDatabaseOutput:     MockOutput,
            ExpectedBody:           "{\"id\":\"id_test\",\"deviceModel\":\"deviceModel_test\",\"name\":\"name_test\",\"note\":\"note_test\",\"serial\":\"serial_test\"}",
            ExpectedStatusCode:     201,
        },
    }

    for _, test := range TestCases {
        // Executing each test cases scenario.
        response := ValidateDatabaseResult(&test.MockDatabaseOutput, test.Error)

        if response.StatusCode != test.ExpectedStatusCode ||  response.Body != test.ExpectedBody{
            t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
        }
    }
}
