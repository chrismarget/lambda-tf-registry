package v1responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/src/common"
)

var _ json.Marshaler = new(Versions)

type Versions struct {
	ItemMaps      []map[string]*dynamodb.AttributeValue
	NamespaceType string
	// Versions      map[string]v1Version
}

type versions struct {
	NamespaceType string    `json:"id"`
	Versions      []version `json:"versions"`
}

type version struct {
	Version   string          `json:"version"`
	Protocols json.RawMessage `json:"Protocols"`
	Platforms []platform      `json:"platforms"`
}

type platform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

func (o Versions) MarshalJSON() ([]byte, error) {
	var err error

	data := versions{NamespaceType: o.NamespaceType}

	versionMap := make(map[string]version)
	for _, itemMap := range o.ItemMaps {
		var versionOsArch string
		err = fetchMapItem(itemMap, "VersionOsArch", &versionOsArch)
		if err != nil {
			return nil, err
		}

		var protocols string
		err = fetchMapItem(itemMap, "Protocols", &protocols)
		if err != nil {
			return nil, err
		}

		parts := strings.Split(versionOsArch, common.PathSep)
		if len(parts) != 3 {
			return nil, errors.New("cannot parse voa from database")
		}

		v := versionMap[parts[0]]
		if string(v.Protocols) != protocols && len(v.Protocols) != 0 {
			return nil, fmt.Errorf(
				"provider %q version %q has conflicting protocol definitions", o.NamespaceType, parts[0],
			)
		}

		v.Version = parts[0]
		v.Protocols = json.RawMessage(protocols)
		v.Platforms = append(v.Platforms, platform{
			Os:   parts[1],
			Arch: parts[2],
		})

		versionMap[parts[0]] = v
	}

	data.Versions = make([]version, len(versionMap))
	var i int
	for _, v := range versionMap {
		data.Versions[i] = v
		i++
	}

	return json.Marshal(data)
}
