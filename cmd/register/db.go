package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/chrismarget/lambda-tf-registry/common/awsclients"
)

const (
	// for filenames like: terraform-provider-jtaf660204634d3ca2a0_0.2.1_SHA256SUMS
	reHashFileParser = "^terraform-provider-([^_]+)_([0-9]+.[0-9]+.[0-9]+)_SHA256SUMS$"

	// for filenames like: terraform-provider-jtaf6601c5cb4d3ca2a0_0.2.1_darwin_amd64.zip
	reZipFileParser1 = "^terraform-provider-([^_]+)_([0-9]+.[0-9]+.[0-9]+)_([^_]+)_([^_]+).zip$"
	reZipFileParser2 = "^(terraform-provider-[^_]+_[0-9]+.[0-9]+.[0-9]+_)[^_]+_[^_]+.zip$"
)

type record struct {
	NamespaceType string `dynamodbav:"NamespaceType"`
	VersionOsArch string `dynamodbav:"VersionOsArch"`
	Keys          string `dynamodbav:"Keys"`
	Protocols     string `dynamodbav:"Protocols"`
	Sha           string `dynamodbav:"SHA"`
	ShaUrl        string `dynamodbav:"SHA_URL"`
	SigUrl        string `dynamodbav:"Sig_URL"`
	Url           string `dynamodbav:"URL"`
}

type updateDbInput struct {
	bucketName            string
	namespace             string
	keysBytes             json.RawMessage
	protocolsBytes        json.RawMessage
	hashFileName          string
	fileHashesToFilePaths map[string]string
	tableName             string
}

func updateDb(ctx context.Context, in updateDbInput) error {
	client, err := awsclients.DdbClient(ctx)
	if err != nil {
		return err
	}

	for fileHash, filePath := range in.fileHashesToFilePaths {
		fileBase := path.Base(filePath)

		re1 := regexp.MustCompile(reZipFileParser1)
		s1 := re1.FindStringSubmatch(fileBase)
		if len(s1) != 5 {
			return fmt.Errorf("failed to parse filename %q with regexp %q",
				path.Base(filePath), reZipFileParser1)
		}

		re2 := regexp.MustCompile(reZipFileParser2)
		s2 := re2.FindStringSubmatch(fileBase)
		if len(s2) != 2 {
			return fmt.Errorf("failed to parse filename %q with regexp %q",
				path.Base(filePath), reZipFileParser2)
		}

		urlbase := fmt.Sprintf("https://%s.s3.amazonaws.com/", in.bucketName)

		r := record{
			NamespaceType: fmt.Sprintf("%s/%s", in.namespace, s1[1]),
			VersionOsArch: path.Join(s1[2:]...),
			Keys:          string(in.keysBytes),
			Protocols:     string(in.protocolsBytes),
			Sha:           fileHash,
			ShaUrl:        urlbase + s2[1] + "SHA256SUMS",
			SigUrl:        urlbase + s2[1] + "SHA256SUMS.sig",
			Url:           urlbase + fileBase,
		}

		item, err := attributevalue.MarshalMap(r)
		if err != nil {
			return fmt.Errorf("failed to marshal new record for database - %w", err)
		}

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: &in.tableName,
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed sending new record to database - %w", err)
		}

		fmt.Printf("%s added to registry database\n", fileBase)
	}

	return nil
}
