package v1handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/chrismarget/lambda-tf-registry/src/common"
)

var _ json.Marshaler = new(ProviderDownloadModel)

type ProviderDownloadModel struct {
	NamespaceType string `dynamodbav:"NamespaceType"`
	VersionOsArch string `dynamodbav:"VersionOsArch"`
	Keys          string `dynamodbav:"Keys"`
	Protocols     string `dynamodbav:"Protocols"`
	Sha           string `dynamodbav:"SHA"`
	ShaUrl        string `dynamodbav:"SHA_URL"`
	SigUrl        string `dynamodbav:"Sig_URL"`
	Url           string `dynamodbav:"URL"`
}

func (o ProviderDownloadModel) GetKey() (map[string]types.AttributeValue, error) {
	namespaceType, err := attributevalue.Marshal(o.NamespaceType)
	if err != nil {
		return nil, err
	}
	versionOsArch, err := attributevalue.Marshal(o.VersionOsArch)
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{
		"NamespaceType": namespaceType,
		"VersionOsArch": versionOsArch,
	}, nil
}

func (o ProviderDownloadModel) MarshalJSON() ([]byte, error) {
	voaParts := strings.Split(o.VersionOsArch, common.PathSep)
	if len(voaParts) != 3 {
		return nil, fmt.Errorf("VersionOsArch should have 3 parts, got: %q", o.VersionOsArch)
	}

	return json.Marshal(&struct {
		Arch      string          `json:"arch"`
		Filename  string          `json:"filename"`
		Keys      json.RawMessage `json:"signing_keys"`
		OS        string          `json:"os"`
		Protocols json.RawMessage `json:"protocols"`
		Sha       string          `json:"shasum"`
		ShaUrl    string          `json:"shasums_url"`
		SigUrl    string          `json:"shasums_signature_url"`
		Url       string          `json:"download_url"`
	}{
		Arch:      voaParts[2],
		Filename:  o.Url[strings.LastIndex(o.Url, common.PathSep)+1:],
		Keys:      []byte(o.Keys),
		OS:        voaParts[1],
		Protocols: []byte(o.Protocols),
		Sha:       o.Sha,
		ShaUrl:    o.ShaUrl,
		SigUrl:    o.SigUrl,
		Url:       o.Url,
	})
}

func NewDownloadModelFromUrlPath(path string) (ProviderDownloadModel, error) {
	var response ProviderDownloadModel

	urlParts := strings.Split(strings.TrimLeft(path, common.PathSep), common.PathSep)
	if len(urlParts) != 8 {
		return response, fmt.Errorf("expected URL to have 8 parts, got %q", path)
	}

	response.NamespaceType = strings.Join(urlParts[2:4], common.PathSep)
	response.VersionOsArch = strings.Join(append(urlParts[4:5], urlParts[6:8]...), common.PathSep)

	return response, nil
}
