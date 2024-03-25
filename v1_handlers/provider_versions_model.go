package v1handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/chrismarget/lambda-tf-registry/common"
)

type ProviderVersionModel struct {
	VersionOsArch string `dynamodbav:"VersionOsArch"`
	Keys          string `dynamodbav:"Keys"`
	Protocols     string `dynamodb:"Protocols"`
}

var _ json.Marshaler = new(ProviderVersionsModel)

type ProviderVersionsModel struct {
	NamespaceType string `dynamodbav:"NamespaceType"`
	Versions      []ProviderVersionModel
}

func (o ProviderVersionsModel) KeyExpr() (expression.Expression, error) {
	keyEx := expression.Key("NamespaceType").Equal(expression.Value(o.NamespaceType))
	return expression.NewBuilder().WithKeyCondition(keyEx).Build()
}

func (o ProviderVersionsModel) MarshalJSON() ([]byte, error) {
	type platform struct {
		Arch string `json:"arch"`
		Os   string `json:"os"`
	}

	type version struct {
		Version   string          `json:"version"`
		Protocols json.RawMessage `json:"protocols"`
		Platforms []platform      `json:"platforms"`
	}

	resultMap := make(map[string]version) // keyed by semantic version string

	for _, v := range o.Versions {
		voaParts := strings.Split(v.VersionOsArch, common.PathSep)
		if len(voaParts) != 3 {
			return nil, fmt.Errorf("VersionOsArch should have 3 parts, got: %q", v.VersionOsArch)
		}

		resultMapItem, ok := resultMap[voaParts[0]]
		if ok && string(resultMapItem.Protocols) != v.Protocols {
			return nil, fmt.Errorf("database has conflicting Protocols entries for %q v%s",
				o.NamespaceType, voaParts[0])
		}

		resultMapItem.Version = voaParts[0]
		resultMapItem.Protocols = json.RawMessage(v.Protocols)
		resultMapItem.Platforms = append(resultMapItem.Platforms, platform{
			Arch: voaParts[2],
			Os:   voaParts[1],
		})

		resultMap[voaParts[0]] = resultMapItem
	}

	result := struct {
		NamespaceType string    `json:"id"`
		Versions      []version `json:"versions"`
	}{
		NamespaceType: o.NamespaceType,
		Versions:      make([]version, len(resultMap)),
	}

	var i int
	for _, v := range resultMap {
		result.Versions[i] = v
		i++
	}

	return json.Marshal(result)
}

func NewVersionsModelFromUrlPath(path string) (ProviderVersionsModel, error) {
	var response ProviderVersionsModel

	urlParts := strings.Split(strings.TrimLeft(path, common.PathSep), common.PathSep)
	if len(urlParts) != 5 {
		return response, fmt.Errorf("expected URL to have 5 parts, got %q", path)
	}

	response.NamespaceType = strings.Join(urlParts[2:4], common.PathSep)

	return response, nil
}
