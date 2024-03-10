// main package
package main

// import other packages
import (
	"os"

	"github.com/Sourjaya/player-api/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Declare a variable of interface type dynamodbiface.
// DynamoDBAPI that enables to mock the DynamoDB service client.

var (
	client dynamodbiface.DynamoDBAPI
)

func main() {
	// get the aws region from the environment variable
	region := os.Getenv("AWS_REGION")

	// Create a new Session passing region as configuration
	Session, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})
	if err != nil {
		return
	}

	// Create a new instance of DynamoDB client
	client = dynamodb.New(Session)

	// Start takes a handler and talks to an internal Lambda endpoint to pass requests to the handler.
	lambda.Start(handler)
}

// DynamoDB table
const tableName = "players"

// Handler function
func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	// if method type is GET
	case "GET":
		return handlers.GetPlayer(req, tableName, client)
	case "POST":
		// if method type is POST
		return handlers.CreatePlayer(req, tableName, client)
	default:
		// for all other method types
		return handlers.Unhandled()
	}
}
