# Simple Serverless Restful API with Golang on AWS Lambda, API Gateway & DynamoDB
## Description
This RESTful API is written on Go language, based on these stacks:
- AWS Lambda
- AWS API Gateway
- AWS DynamoDB
- Serverless Framework
## API Request-Responce Cycle
The API accepts the following JSON requests and produces the corresponding HTTP responses:
### Request 1:
Request to insert a new device to database(DynamoDB).
```
HTTP Method: POST
URL: https://<api-gateway-url>/api/devices
content-type: application/json
Body:
  {
    "id": "/devices/id1",
    "deviceModel": "/devicemodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
  }
```
#### Response 1 - Success:
Provided data inserted to database(DynamoDB) successfully.
```
HTTP-Statuscode: HTTP 201
content-type: application/json
Body:
  {
    "id": "/devices/id1",
    "deviceModel": "/devicemodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
  }
```
#### Response 1 - Failure 1:
If any of the payload fields are missing, response will have a descriptive error message for client.
```
HTTP-Statuscode: HTTP 400
"Following fields are not provided: id, serial, ..."
```
#### Response 1 - Failure 2:
If any exceptional situation occurs on the server side.

```
HTTP-Statuscode: HTTP 500
"Internal Server's Error occurred."
```
### Request 2:
Get a device based on provided id.
```
HTTP Method: GET
URL: https://<api-gateway-url>/api/devices/{id}

Replace {id} with desire device id
```
#### Response 2 - Success:
The desire id exists on DynamoDB.
```
HTTP-Statuscode: HTTP 200
content-type: application/json
body:
  {
    "id": "/devices/id1",
    "deviceModel": "/devicemodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
  }
```
#### Response 2 - Failure 1:
```
HTTP-Statuscode: HTTP 404
"Desired device with provided id was not founded."
```
#### Response 2 - Failure 2:
If any exceptional situation occurs on the server side.
```
HTTP-Statuscode: HTTP 500
"Internal Server's Error occured."
```
## API Included:
- [`script`](https://github.com/parhizi/simple-go-restful-aws/tree/master/scripts) folder contains three bash script files which automate the process of build, depoly and test.
- [`addDevice.go`](https://github.com/parhizi/simple-go-restful-aws/blob/master/src/handlers/addDevice/addDevice.go) is responsible for adding desire items to the DynamoDB based on the database schema.
- [`getDeviceById.go`](https://github.com/parhizi/simple-go-restful-aws/blob/master/src/handlers/getDeviceById/getDeviceById.go) is responsible for making query based on the given id.
- [`addDevice_test.go`](https://github.com/parhizi/simple-go-restful-aws/blob/master/src/handlers/addDevice/addDevice_test.go) and [`getDeviceById_test.go`](https://github.com/parhizi/simple-go-restful-aws/blob/master/src/handlers/getDeviceById/getDeviceById_test.go) contain all the test case scenarios.
- [`serverless.yml`](https://github.com/parhizi/simple-go-restful-aws/blob/master/serverless.yml) have Serverless Framework configurations which will set AWS services on behalf of you.
## Dependencies
For deploying this API, you need to install and configure the following items:
- [`Go`](https://golang.org/) Because this API is written on it! :)
- [`Serverless Framework`](https://serverless.com/): Automating deployment on AWS.
- [`dep`](https://golang.github.io/dep/): Dependency management tool for Go.
- `Bash` In case of running build, deploy and test script files.
Note that you have to configure all above items based on your operating system variables and AWS services.
## Get the things work
After satisfying above prerequisites, clone the project in your desire folder and and run these scripts within that folder, based on the following needs:
### Building
First, this script run `dep init` and `dep ensure` to provide the packages needed for building process. After that it `go build` all API files.
```
./script/build.sh
```
### Deploying
This script will deploy the API based on the `serverless.yml` configuration file to the AWS.
```
./script/deploy.sh
```
### Unit Testing
By executing this script, `*_test.go` file of each `addDevice.go` and `getDeviceById.go` will be executed. At last the script will save the test coverage result in `cover.html` file in each of the function's folder.
```
./script/build.sh
```
#### Unit Test Output Sample 
```
Testing .go files
PASS
coverage: 91.8% of statements
ok  	github.com/me/Simple-Go-RESTful-AWS/src/handlers/addDevice	0.026s
PASS
coverage: 93.1% of statements
ok  	github.com/me/Simple-Go-RESTful-AWS/src/handlers/getDeviceById	0.005s
Done.
```
## Testing in real world:
We can have real world testing with AWS endpoints, provided to us after deploying the API to AWS. We test our both HTTP global verbs by [`cURL`](https://curl.haxx.se/), a command line tool and library for transferring data with URLs.
### PUT sample:
```
curl -i -H "Content-Type: application/json" -X POST https://<api-gateway-url>/devices -d '{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}'

Response:
HTTP-Statuscode: HTTP 201
content-type: application/json
Body:
  {
    "id": "/devices/id1",
    "deviceModel": "/devicemodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
  }
```
### GET sample:
Let's query the previously added item.
```
curl -i https://<api-gateway-url>/devices/id1HTTP-Statuscode: HTTP 200

Response:
content-type: application/json
body:
  {
    "id": "/devices/id1",
    "deviceModel": "/devicemodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
  }
```
