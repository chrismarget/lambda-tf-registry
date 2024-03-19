package responders

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"net/http"
	"strings"
)

var _ json.Marshaler = new(V1Versions)
var _ Responder = new(V1Versions)

type V1Versions struct {
	Items         []map[string]*dynamodb.AttributeValue
	NamespaceType string
	Versions      map[string]v1Version
}

type v1VersionPlatform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

type v1Version struct {
	Version   string              `json:"version"`
	Protocols json.RawMessage     `json:"protocols"`
	Platforms []v1VersionPlatform `json:"platforms"`
}

//func (o *v1Version) LoadFromItemMap(itemMap map[string]*dynamodb.AttributeValue) error {
//	var voa string
//	err := fetchMapItem(itemMap, "VersionOsArch", &voa)
//	if err != nil {
//		return err
//	}
//
//	voaSlice := strings.Split(voa, "/")
//	if len(voaSlice) != 3 {
//		return fmt.Errorf("VersionOsArch has unexpected parts count (%d): %q", len(voaSlice), voa)
//	}
//
//	o.Version = voaSlice[0]
//
//	err = fetchMapItem(itemMap, "Protocols", &o.Protocols)
//	if err != nil {
//		return err
//	}
//
//	for o.Platforms =
//
//}

func (o *V1Versions) MarshalJSON() ([]byte, error) {

	type platform struct {
		Os   string `json:"os"`
		Arch string `json:"arch"`
	}

	type version struct {
		Version   string          `json:"version"`
		Protocols json.RawMessage `json:"protocols"`
		Platforms []platform      `json:"platforms"`
	}

	renderVersion := func(itemMap map[string]*dynamodb.AttributeValue) (version, error) {
		var err error
		var result version

		var voa string
		err = fetchMapItem(itemMap, "VersionOsArch", &voa)
		if err != nil {
			return result, err
		}

		voaSlice := strings.Split(voa, "/")
		if len(voaSlice) != 3 {
			return result, fmt.Errorf("VersionOsArch has unexpected parts count (%d): %q", len(voaSlice), voa)
		}

		result.Version = voaSlice[0]

		err = fetchMapItem(itemMap, "Protocols", &result.Protocols)

		return result, nil
	}

	s := struct {
		NamespaceType string    `json:"id"`
		Versions      []version `json:"versions"`
	}{
		NamespaceType: o.NamespaceType,
		Versions:      make([]version, len(o.Items)),
	}

	for i := range o.Items {
		v, err := renderVersion(o.Items[i])
		if err != nil {
			return nil, err
		}
		s.Versions[i] = v
	}

	return json.Marshal(s)
}

func (o *V1Versions) AddItem(itemMap map[string]*dynamodb.AttributeValue) error {
	var voa string
	err := fetchMapItem(itemMap, "VersionOsArch", &voa)
	if err != nil {
		return fmt.Errorf("failed fetching VersionOsArch - %w", err)
	}

	voaSlice := strings.Split(voa, "/")
	if len(voaSlice) != 3 {
		return fmt.Errorf("VersionOsArch should have 3 parts, got %q", voa)
	}

	version := o.Versions[voaSlice[0]]
	version.Version = voaSlice[0]
	err = fetchMapItem(itemMap, "Protocols", &version.Protocols)
	if err != nil {
		return fmt.Errorf("failed fetching Protocols - %w", err)
	}

	version.Platforms = append(version.Platforms, v1VersionPlatform{
		Os:   voaSlice[1],
		Arch: voaSlice[2],
	})

	o.Versions[voaSlice[0]] = version

	return nil
}

func (o *V1Versions) Respond() (events.LambdaFunctionURLResponse, error) {
	if o.Versions == nil {
		o.Versions = make(map[string]v1Version)
	}

	for _, item := range o.Items {
		err := o.AddItem(item)
		if err != nil {
			return Error{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("failed adding item - %w", err),
			}.Respond()
		}
	}

	s := struct {
		NamespaceType string      `json:"id"`
		Versions      []v1Version `json:"versions"`
	}{
		NamespaceType: o.NamespaceType,
		Versions:      make([]v1Version, len(o.Versions)),
	}

	var i int
	for _, v := range o.Versions {
		s.Versions[i] = v
		i++
	}

	body, err := json.Marshal(s)
	if err != nil {
		return Error{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed marshaling JSON - %w", err),
		}.Respond()
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
