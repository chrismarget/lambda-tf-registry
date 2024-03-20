package responders

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrismarget/lambda-tf-registry/src/env"
	"github.com/chrismarget/lambda-tf-registry/src/utils"
	"log"
	"net/http"
	"path"
)

var _ Responder = new(RegisterProvider)

type RegisterProvider struct {
	Arch                string          `json:"arch"`
	DownloadUrl         string          `json:"download_url"`
	Os                  string          `json:"os"`
	Protocols           json.RawMessage `json:"protocols"`
	SigningKeys         json.RawMessage `json:"signing_keys"`
	Shasum              string          `json:"shasum"`
	ShasumsUrl          string          `json:"shasums_url"`
	ShasumsSignatureUrl string          `json:"shasums_signature_url"`
	Version             string          `json:"version"`
	Namespace           string          `json:"namespace"`
	Type                string          `json:"type"`
}

func (o RegisterProvider) Respond() (events.LambdaFunctionURLResponse, error) {
	log.Printf("registerProvider: %s", o)
	//body, err := json.Marshal(o)
	//if err != nil {
	//	return Error{
	//		Code: http.StatusInternalServerError,
	//		Err:  err,
	//	}.Respond()
	//}

	err := o.loadRecord()
	if err != nil {
		return Error{
			Code: http.StatusInternalServerError,
			Err:  err,
		}.Respond()
	}

	response := events.LambdaFunctionURLResponse{StatusCode: http.StatusOK}

	return response, nil
}

func (o RegisterProvider) loadRecord() error {
	client, err := utils.DdbClient()
	if err != nil {
		return fmt.Errorf("failed getting dynamodb client - %w", err)
	}

	tableName := env.Get(env.ProviderTableName)
	if tableName == "" {
		return fmt.Errorf("env var %q not set", env.ProviderTableName)
	}

	item := struct {
		NamespaceType string `dynamodbav:"NamespaceType"`
		VersionOsArch string `dynamodbav:"VersionOsArch"`
		Keys          string `dynamodbav:"Keys"`
		Protocols     string `dynamodbav:"Protocols"`
		SHA           string `dynamodbav:"SHA"`
		ShaUrl        string `dynamodbav:"SHA_URL"`
		SigUrl        string `dynamodbav:"Sig_URL"`
		Url           string `dynamodbav:"URL"`
	}{
		NamespaceType: path.Join(o.Namespace, o.Type),
		VersionOsArch: path.Join(o.Version, o.Os, o.Arch),
		Keys:          string(o.SigningKeys),
		Protocols:     string(o.Protocols),
		SHA:           o.Shasum,
		ShaUrl:        o.ShasumsSignatureUrl,
		SigUrl:        o.ShasumsSignatureUrl,
		Url:           o.DownloadUrl,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed marshaling dynamodb attributes - %w", err)
	}

	_, err = client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &tableName,
	})
	if err != nil {
		return fmt.Errorf("failed putting dynamodb item - %w", err)
	}

	return nil
}
