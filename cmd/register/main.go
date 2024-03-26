package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
)

const usage = `

To use this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

  terraform {
    required_providers {
      jtaf = {
        source = %q 
        version = %q
      }
    }
  }

  provider "jtaf" {
    # Configuration options
  }
`

var (
	flagDir       string
	flagProtocols string
	flagBucket    string
	flagNamespace string
	flagRegistry  string
	flagTable     string

	providerType    string
	providerVersion string
)

func init() {
	flag.StringVar(&flagProtocols, "protocols", `["5.0"]`, "list of terraform plugin protocols")
	flag.StringVar(&flagDir, "dir", "./dist", "release directory")
	flag.StringVar(&flagBucket, "bucket", "jtaf-registry", "s3 bucket name")
	flag.StringVar(&flagNamespace, "namespace", "juniper", "namespace (publisher) on the registry")
	flag.StringVar(&flagRegistry, "registry", "tf-registry.click", "registry hostname")
	flag.StringVar(&flagTable, "table", "registry-providers", "dynamodb table name")
}

func parseFlags() {
	flag.Parse()

	err := typeAndVersion(flagDir)
	if err != nil {
		log.Fatalf("failed trying to determine provider type and/or version - %s", err)
	}

	fi, err := os.Stat(flagDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("directory %q not found", flagDir)
		}
	}
	if !fi.IsDir() {
		log.Fatalf("%q is not a directory", flagDir)
	}
}

func main() {
	ctx := context.Background()

	parseFlags()

	hashFileName, err := getHashFile(flagDir)
	if err != nil {
		log.Fatal(err)
	}

	fileHashesToFilePaths, err := hashes(hashFileName)
	if err != nil {
		log.Fatal(err)
	}

	keysBytes, err := getKeys(flagDir)
	if err != nil {
		log.Fatal(err)
	}

	err = uploadFiles(ctx, uploadFilesInput{
		bucketName:            flagBucket,
		protocolsBytes:        []byte(flagProtocols),
		hashFileName:          hashFileName,
		fileHashesToFilePaths: fileHashesToFilePaths,
		keysBytes:             keysBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = updateDb(ctx, updateDbInput{
		bucketName:            flagBucket,
		namespace:             flagNamespace,
		protocolsBytes:        []byte(flagProtocols),
		hashFileName:          hashFileName,
		fileHashesToFilePaths: fileHashesToFilePaths,
		keysBytes:             keysBytes,
		tableName:             flagTable,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(usage, path.Join(flagRegistry, "juniper", providerType), providerVersion)

	os.Exit(0)
}

func typeAndVersion(dir string) error {
	if providerType != "" && providerVersion != "" {
		return nil
	}

	hashFileName, err := getHashFile(dir)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(reHashFileParser)
	s := re.FindStringSubmatch(path.Base(hashFileName))
	if len(s) != 3 {
		return fmt.Errorf("failed to parse filename %q with regexp %q",
			path.Base(hashFileName), reHashFileParser)
	}

	providerType = s[1]
	providerVersion = s[2]

	return nil
}
