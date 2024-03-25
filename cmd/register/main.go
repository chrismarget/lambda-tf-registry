package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
)

const envVarProviderType = "PTYPE"

var (
	flagDist      string
	flagProtocols string
	flagBucket    string
	flagNamespace string
	// flagType      string
	flagTable string
)

func init() {
	flag.StringVar(&flagProtocols, "protocols", `["5.0"]`, "list of terraform plugin protocols")
	flag.StringVar(&flagDist, "dist", "./dist", "release directory")
	flag.StringVar(&flagBucket, "bucket", "jtaf-registry", "s3 bucket name")
	flag.StringVar(&flagNamespace, "namespace", "juniper", "namespace (publisher) on the registry")
	// flag.StringVar(&flagType, "type", os.Getenv(envVarProviderType), "type (provider name) on the registry")
	flag.StringVar(&flagTable, "table", "registry-providers", "dynamodb table name")
}

func parseFlags() {
	flag.Parse()

	fi, err := os.Stat(flagDist)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("directory %q not found", flagDist)
		}
	}
	if !fi.IsDir() {
		log.Fatalf("%q is not a directory", flagDist)
	}
}

func main() {
	ctx := context.Background()

	parseFlags()

	hashFileName, err := getHashFile(flagDist)
	if err != nil {
		log.Fatal(err)
	}

	fileHashesToFilePaths, err := hashes(hashFileName)
	if err != nil {
		log.Fatal(err)
	}

	keysBytes, err := getKeys(flagDist)
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

	os.Exit(0)
}
