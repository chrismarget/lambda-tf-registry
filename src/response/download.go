package response

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

var _ json.Marshaler = new(Download)

type Download struct {
	NamespaceType string
	VersionOsArch string

	keys      string
	protocols string
	sha       string
	shaUrl    string
	sigUrl    string
	url       string
}

func (o *Download) MarshalJSON() ([]byte, error) {
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

func (o *Download) Respond(itemMap map[string]*dynamodb.AttributeValue) (events.LambdaFunctionURLResponse, error) {
	if o.NamespaceType == "" {
		return events.LambdaFunctionURLResponse{}, errors.New("o.NamespaceType must not be empty")
	}

	if o.VersionOsArch == "" {
		return events.LambdaFunctionURLResponse{}, errors.New("o.VersionOsArch must not be empty")
	}

	if itemMap == nil {
		return events.LambdaFunctionURLResponse{}, errors.New("itemMap must not be nil")
	}

	err := o.loadItems(itemMap)
	if err != nil {
		return ierrors.IErr{
			Err:  fmt.Errorf("failed while loading items - %w", err),
			Code: http.StatusInternalServerError,
		}.LambdaResponse()
	}

	data, err := json.Marshal(o)
	if err != nil {
		return ierrors.IErr{
			Err:  fmt.Errorf("failed while marshaling response - %w", err),
			Code: http.StatusInternalServerError,
		}.LambdaResponse()
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(data),
	}, nil
}

func FetchMapItem(itemMap map[string]*dynamodb.AttributeValue, item string, target any) error {
	v, ok := itemMap[item]
	if !ok {
		return fmt.Errorf("item %q not found in map", item)
	}

	switch t := target.(type) {
	case *string:
		*t = *v.S
	default:
		return fmt.Errorf("unhandled type: %T", t)
	}

	return nil
}

func (o *Download) loadItems(itemMap map[string]*dynamodb.AttributeValue) error {
	var err error

	err = FetchMapItem(itemMap, "Keys", &o.keys)
	if err != nil {
		return err
	}

	err = FetchMapItem(itemMap, "Protocols", &o.protocols)
	if err != nil {
		return err
	}

	err = FetchMapItem(itemMap, "SHA", &o.sha)
	if err != nil {
		return err
	}

	err = FetchMapItem(itemMap, "SHA_URL", &o.shaUrl)
	if err != nil {
		return err
	}

	err = FetchMapItem(itemMap, "Sig_URL", &o.sigUrl)
	if err != nil {
		return err
	}

	err = FetchMapItem(itemMap, "URL", &o.url)
	if err != nil {
		return err
	}

	return nil
}
