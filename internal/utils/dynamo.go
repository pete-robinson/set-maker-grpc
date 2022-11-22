package utils

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)


func CreateDynamoClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}


// Builds a base64 encoded byte slice from a dynamodb attribute value map
// Used for pagination cursors
func EncodeAttributeMap(in map[string]types.AttributeValue) (string, error) {
	if in == nil {
		return "", nil
	}

	// convert attrmap to go values
	var gov map[string]interface{}
	if err := attributevalue.UnmarshalMap(in, &gov); err != nil {
		return "", err
	}

	// marshal map to []bytes
	jsn, err := json.Marshal(gov)
	if err != nil {
		return "", err
	}

	// base64 encode []bytes to string
	return base64.StdEncoding.EncodeToString(jsn), nil
}


func DecodeAttributeMap(in string) (map[string]types.AttributeValue, error) {
	if in == "" {
		return nil, nil
	}

	// base64 decode
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}

	// now we unmarshal the json
	var key map[string]string
	if err = json.Unmarshal(decoded, &key); err != nil {
		return nil, err
	}

	// finally pass the resulting map to Dynamo's helpers to create a map[string]types.AttributeValue
	return attributevalue.MarshalMap(key)
}
