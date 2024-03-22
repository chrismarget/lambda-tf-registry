package v1responses

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func fetchMapItem(itemMap map[string]*dynamodb.AttributeValue, item string, target any) error {
	v, ok := itemMap[item]
	if !ok {
		return fmt.Errorf("item %q not found in map", item)
	}

	switch t := target.(type) {
	case *string:
		*t = *v.S
	case *json.RawMessage:
		*t = json.RawMessage(*v.S)
	default:
		return fmt.Errorf("unhandled type: %T", t)
	}

	return nil
}
