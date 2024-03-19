package responders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ierrors "github.com/chrismarget/lambda-tf-registry/src/errors"
)

var _ json.Marshaler = new(V1Download)
var _ Responder = new(V1Download)

type V1Download struct {
	ItemMap       map[string]*dynamodb.AttributeValue
	NamespaceType string
	VersionOsArch string

	keys      string
	protocols string
	sha       string
	shaUrl    string
	sigUrl    string
	url       string
}

func (o *V1Download) MarshalJSON() ([]byte, error) {
	voaParts := strings.Split(o.VersionOsArch, "/")
	if len(voaParts) != 3 {
		return nil, fmt.Errorf("VersionOsArch should have 3 parts, got: %q", o.VersionOsArch)
	}

	s := struct {
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
		Filename:  o.url[strings.LastIndex(o.url, "/")+1:],
		Keys:      []byte(o.keys),
		OS:        voaParts[1],
		Protocols: []byte(o.protocols),
		Sha:       o.sha,
		ShaUrl:    o.shaUrl,
		SigUrl:    o.sigUrl,
		Url:       o.url,
	}

	return json.Marshal(&s)
}

func (o *V1Download) Respond() (events.LambdaFunctionURLResponse, error) {
	if o.ItemMap == nil {
		return events.LambdaFunctionURLResponse{}, errors.New("o.ItemMap must not be nil")
	}

	err := o.loadItems(o.ItemMap)
	if err != nil {
		return ierrors.IErr{
			Err:  fmt.Errorf("failed while loading items - %w", err),
			Code: http.StatusInternalServerError,
		}.LambdaResponse()
	}

	data, err := json.Marshal(o)
	if err != nil {
		return ierrors.IErr{
			Err:  fmt.Errorf("failed while marshaling responders - %w", err),
			Code: http.StatusInternalServerError,
		}.LambdaResponse()
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(data),
	}, nil
}

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

func (o *V1Download) loadItems(itemMap map[string]*dynamodb.AttributeValue) error {
	var err error

	err = fetchMapItem(itemMap, "Keys", &o.keys)
	if err != nil {
		return err
	}

	err = fetchMapItem(itemMap, "Protocols", &o.protocols)
	if err != nil {
		return err
	}

	err = fetchMapItem(itemMap, "SHA", &o.sha)
	if err != nil {
		return err
	}

	err = fetchMapItem(itemMap, "SHA_URL", &o.shaUrl)
	if err != nil {
		return err
	}

	err = fetchMapItem(itemMap, "Sig_URL", &o.sigUrl)
	if err != nil {
		return err
	}

	err = fetchMapItem(itemMap, "URL", &o.url)
	if err != nil {
		return err
	}

	return nil
}
