package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const envVarProviderType = "PTYPE"

var (
	distDir    string
	protocols  string
	bucketName string
	namespace  string
	ptype      string
)

func init() {
	flag.StringVar(&protocols, "protocols", `["5.0"]`, "list of terraform plugin protocols")
	flag.StringVar(&distDir, "dist", "./dist", "release directory")
	flag.StringVar(&bucketName, "bucket", "jtaf-registry", "s3 bucket name")
	flag.StringVar(&namespace, "namespace", "juniper", "namespace (publisher) on the registry")
	flag.StringVar(&ptype, "type", os.Getenv(envVarProviderType), "type (provider name) on the registry")
	flag.Parse()

	if ptype == "" {
		fmt.Print("\nError: --ptype argument is required.\n\n")
		flag.Usage()
		os.Exit(1)
	}
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

	payload := updateDbInput{
		namespaceType:  fmt.Sprintf("%s/%s", namespace, ptype),
		protocolsBytes: []byte(protocols),
	}

	var err error

	payload.keysBytes, err = getKeys(distDir)
	if err != nil {
		log.Fatal(err)
	}

	payload.hashFileName, err = getHashFile(distDir)
	if err != nil {
		log.Fatal(err)
	}

	err = updateDb(ctx, payload)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

//func uploadSharedObjects(ctx context.Context, hashfile string) error {
//	s3UploadManager, err := awsclients.S3Manager(ctx)
//	if err != nil {
//		return err
//	}
//
//	fHash, err := os.Open(hashfile)
//	if err != nil {
//		return fmt.Errorf("failed to open %q for S3 upload - %w", hashfile, err)
//	}
//	defer func() {
//		fHash.Close()
//	}()
//
//	_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
//		Bucket: aws.String(bucketName),
//		Key:    aws.String(path.Base(hashfile)),
//		Body:   fHash,
//	})
//	if err != nil {
//		return fmt.Errorf("failed while uploading %q to s3://%s - %w", hashfile, bucketName, err)
//	}
//
//	fSig, err := os.Open(hashfile + sigFileSuffix)
//	if err != nil {
//		return fmt.Errorf("failed to open %q for S3 upload - %w", hashfile+sigFileSuffix, err)
//	}
//	defer func() {
//		fHash.Close()
//	}()
//
//	_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
//		Bucket: aws.String(bucketName),
//		Key:    aws.String(path.Base(hashfile + sigFileSuffix)),
//		Body:   fSig,
//	})
//	if err != nil {
//		return fmt.Errorf("failed while uploading %q to s3://%s - %w", hashfile+sigFileSuffix, bucketName, err)
//	}
//
//	return nil
//}
//
//func uploadProviderFiles(ctx context.Context, fileNamesToHashes map[string]string) error {
//	s3UploadManager, err := awsclients.S3Manager(ctx)
//	if err != nil {
//		return err
//	}
//
//	for hash, fName := range fileNamesToHashes {
//		f, err := os.Open(fName)
//		if err != nil {
//			return fmt.Errorf("failed to open %q for S3 upload - %w", fName, err)
//		}
//
//		bSum, err := hex.DecodeString(hash)
//		if err != nil {
//			return fmt.Errorf("failed decoding hex string %q - %w", hash, err)
//		}
//
//		b64Sum := base64.StdEncoding.EncodeToString(bSum)
//		_, err = s3UploadManager.Upload(ctx, &s3.PutObjectInput{
//			Bucket:         aws.String(bucketName),
//			Key:            aws.String(path.Base(fName)),
//			Body:           f,
//			ChecksumSHA256: &b64Sum,
//		})
//		if err != nil {
//			return fmt.Errorf("failed while uploading %q to s3://%s - %w", fName, bucketName, err)
//		}
//
//		err = f.Close()
//		if err != nil {
//			return fmt.Errorf("failed closing file %q - %w", fName, err)
//		}
//	}
//
//	return nil
//}
