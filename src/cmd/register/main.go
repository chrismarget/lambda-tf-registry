package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chrismarget/lambda-tf-registry/src/common/awsclients"
	"log"
	"os"
	"path"
)

var (
	distDir    string
	protocols  string
	bucketName string
)

func init() {
	flag.StringVar(&protocols, "protocols", `["5.0"]`, "list of terraform plugin protocols")
	flag.StringVar(&distDir, "dist", "./dist", "release directory")
	flag.StringVar(&bucketName, "bucket", "jtaf-registry", "s3 bucket name")
	flag.Parse()

	checkDir(distDir)
}

func checkDir(d string) {
	fi, err := os.Stat(d)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("directory %q not found", d)
		}
	}
	if !fi.IsDir() {
		log.Fatalf("%q is not a directory", d)
	}
}

func main() {
	ctx := context.Background()

	namespaceType, err := getNamespaceType()
	if err != nil {
		log.Fatal(err)
	}
	_ = namespaceType

	keys, err := getKeys(distDir)
	if err != nil {
		log.Fatal(err)
	}
	_ = keys

	protocolBytes := []byte(protocols)
	_ = protocolBytes

	hashFile, err := getHashFile(distDir)
	if err != nil {
		log.Fatal(err)
	}

	hashesToFiles, err := hashes(hashFile)
	if err != nil {
		log.Fatal(err)
	}
	_ = hashesToFiles

	err = processSharedObjects(ctx, hashFile)
	if err != nil {
		log.Fatal(err)
	}

	err = processProviderFiles(ctx, hashesToFiles)
	if err != nil {
		log.Fatal(err)
	}
}

func processSharedObjects(ctx context.Context, hashfile string) error {
	s3UploadManager, err := awsclients.S3Manager(ctx)
	if err != nil {
		return err
	}

	fHash, err := os.Open(hashfile)
	if err != nil {
		return fmt.Errorf("failed to open %q for S3 upload - %w", hashfile, err)
	}
	defer func() {
		fHash.Close()
	}()

	_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path.Base(hashfile)),
		Body:   fHash,
	})
	if err != nil {
		return fmt.Errorf("failed while uploading %q to s3://%s - %w", hashfile, bucketName, err)
	}

	fSig, err := os.Open(hashfile + sigFileSuffix)
	if err != nil {
		return fmt.Errorf("failed to open %q for S3 upload - %w", hashfile+sigFileSuffix, err)
	}
	defer func() {
		fHash.Close()
	}()

	_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path.Base(hashfile + sigFileSuffix)),
		Body:   fSig,
	})
	if err != nil {
		return fmt.Errorf("failed while uploading %q to s3://%s - %w", hashfile+sigFileSuffix, bucketName, err)
	}

	return nil
}

func processProviderFiles(ctx context.Context, fileNamesToHashes map[string]string) error {
	s3UploadManager, err := awsclients.S3Manager(ctx)
	if err != nil {
		return err
	}

	for hash, fName := range fileNamesToHashes {
		f, err := os.Open(fName)
		if err != nil {
			return fmt.Errorf("failed to open %q for S3 upload - %w", fName, err)
		}

		bSum, err := hex.DecodeString(hash)
		if err != nil {
			return fmt.Errorf("failed decoding hex string %q - %w", hash, err)
		}

		b64Sum := base64.StdEncoding.EncodeToString(bSum)
		_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
			Bucket:         aws.String(bucketName),
			Key:            aws.String(path.Base(fName)),
			Body:           f,
			ChecksumSHA256: &b64Sum,
		})
		if err != nil {
			return fmt.Errorf("failed while uploading %q to s3://%s - %w", fName, bucketName, err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed closing file %q - %w", fName, err)
		}
	}

	return nil
}
