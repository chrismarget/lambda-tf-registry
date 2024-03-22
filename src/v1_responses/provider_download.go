package v1responses

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/src/common"
)

var _ json.Marshaler = new(Download)

type Download struct {
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

type download struct {
	Arch      string          `json:"arch"`
	Filename  string          `json:"filename"`
	Keys      json.RawMessage `json:"signing_keys"`
	OS        string          `json:"os"`
	Protocols json.RawMessage `json:"protocols"`
	Sha       string          `json:"shasum"`
	ShaUrl    string          `json:"shasums_url"`
	SigUrl    string          `json:"shasums_signature_url"`
	Url       string          `json:"download_url"`
}

func (o *Download) loadItems() error {
	var err error

	err = fetchMapItem(o.ItemMap, "Keys", &o.keys)
	if err != nil {
		return err
	}

	err = fetchMapItem(o.ItemMap, "Protocols", &o.protocols)
	if err != nil {
		return err
	}

	err = fetchMapItem(o.ItemMap, "SHA", &o.sha)
	if err != nil {
		return err
	}

	err = fetchMapItem(o.ItemMap, "SHA_URL", &o.shaUrl)
	if err != nil {
		return err
	}

	err = fetchMapItem(o.ItemMap, "Sig_URL", &o.sigUrl)
	if err != nil {
		return err
	}

	err = fetchMapItem(o.ItemMap, "URL", &o.url)
	if err != nil {
		return err
	}

	return nil
}

func (o *Download) MarshalJSON() ([]byte, error) {
	err := o.loadItems()
	if err != nil {
		return nil, err
	}

	voaParts := strings.Split(o.VersionOsArch, "/")
	if len(voaParts) != 3 {
		return nil, fmt.Errorf("VersionOsArch should have 3 parts, got: %q", o.VersionOsArch)
	}

	result := download{
		Arch:      voaParts[2],
		Filename:  o.url[strings.LastIndex(o.url, common.PathSep)+1:],
		Keys:      []byte(o.keys),
		OS:        voaParts[1],
		Protocols: []byte(o.protocols),
		Sha:       o.sha,
		ShaUrl:    o.shaUrl,
		SigUrl:    o.sigUrl,
		Url:       o.url,
	}

	return json.Marshal(&result)
}
