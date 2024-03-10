// This code is part of player package
package player

// import other packages
import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

// different error messages
var (
	ErrorInvalidID               = "Invalid UUID"
	ErrorFailedToUnmarshalRecord = "Failed to unmarshal record"
	ErrorFailedToFetchRecord     = "Failed to fetch record"
	ErrorInvalidPlayerData       = "Invalid Player data"
	ErrorCouldNotMarshalItem     = "Could not marshal item"
	ErrorCouldNotPostItem        = "Could not post item in DB"
	ErrorPlayerAlreadyExists     = "Player already exists"
	ErrorPlayerDoesNotExist      = "Player does not exist"
)

type Player struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Position  string `json:"position"`
	Club      string `json:"club"`
}

// Fetch Player details from DynamoDB using ID
func GetPlayerByID(id, tableName string, client dynamodbiface.DynamoDBAPI) (*Player, error) {

	//regular expression used to check for validity of the ID given as URL parameter.
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	if !r.MatchString(id) {
		return nil, errors.New(ErrorInvalidID)
	}
	// create the input of get item operation, which contains the id as key value to search in the table.
	// We also need to specify the table name.
	// dynamodb.AttributeValue represents data for an attribute
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	// call the GetItem method to fetch data from dynamoDB table.
	result, err := client.GetItem(input)
	// Check if the result is nil or if there is any error during fetching the record.
	if err != nil || result.Item == nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	// Convert the json data into Player data structure
	item := new(Player)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

// Fetch details of all players
func GetPlayers(tableName string, client dynamodbiface.DynamoDBAPI) (*[]Player, error) {
	// create the input of scan operation which contains the table name.
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Scan the DynamoDB table and store the output.
	result, err := client.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Player)
	// Convert the json data into a slice of Player data structure
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

// Store data for a new entry of a player.
func CreatePlayer(req events.APIGatewayProxyRequest, tableName string, client dynamodbiface.DynamoDBAPI) (
	*Player,
	error,
) {
	var p Player

	// Convert the json data received from request body, into Player data structure.
	if err := json.Unmarshal([]byte(req.Body), &p); err != nil {
		return nil, errors.New(ErrorInvalidPlayerData)
	}
	// Create a new UUID and assign to the Player data, as id.
	id := uuid.New()
	p.ID = id.String()

	checkPlayer, _ := GetPlayerByID(p.ID, tableName, client)
	if checkPlayer != nil && len(checkPlayer.ID) != 0 {
		return nil, errors.New(ErrorPlayerAlreadyExists)
	}
	// Convert the data from Player data structure to json data.
	item, err := dynamodbattribute.MarshalMap(p)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	// input data to be put in the DynamoDB table
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	// Check for error in operation
	_, err = client.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotPostItem)
	}

	return &p, nil
}
