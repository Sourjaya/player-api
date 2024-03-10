// This code is part of handlers package
package handlers

// import other packages
import (
	"net/http"

	"github.com/Sourjaya/player-api/pkg/player"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// create a error message variable
const ErrorMethodNotAllowed = "method not allowed"

// Create a struct to hold aws error as string
type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// this function will be called when the method type is "GET"
func GetPlayer(req events.APIGatewayProxyRequest, tableName string, client dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error,
) {
	// if the name is provided as URL parameter then call GetPlayerByName function from player package.
	id := req.QueryStringParameters["id"]

	if id != "" {
		result, err := player.GetPlayerByID(id, tableName, client)
		if err != nil {
			return response(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return response(http.StatusOK, result)
	}
	// else call GetPlayers function from player package.
	result, err := player.GetPlayers(tableName, client)
	if err != nil {
		return response(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return response(http.StatusOK, result)
}

// this function will be called when the method type is "POST"
func CreatePlayer(req events.APIGatewayProxyRequest, tableName string, client dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error,
) {
	// Call CreatePlayer function from player package.
	result, err := player.CreatePlayer(req, tableName, client)
	if err != nil {
		return response(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return response(http.StatusCreated, result)
}

// Function that return a Method not allowed error message for those unhandled method types.
func Unhandled() (*events.APIGatewayProxyResponse, error) {
	return response(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
