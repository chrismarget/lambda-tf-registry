package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chrismarget/lambda-tf-registry/common/awsclients"
)

type updateDbInput struct {
	bucketName        string
	namespaceType     string
	keysBytes         json.RawMessage
	protocolsBytes    json.RawMessage
	hashFileName      string
	hashesToFilenames map[string]string
	wg                *sync.WaitGroup
}

func updateDb(ctx context.Context, in updateDbInput) error {
	// read the file with checksums of each provider binary
	hashesToFilenames, err := hashes(in.hashFileName)
	if err != nil {
		log.Fatal(err)
	}
	_ = hashesToFilenames

	s3UploadManager, err := awsclients.S3Manager(ctx)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error)

	ufcfg := uploadFileIn{
		ctx:     ctx,
		mgr:     s3UploadManager,
		errChan: errChan,
		bucket:  in.bucketName,
		wg:      wg,
	}

	// send the hash file
	ufcfg.localpath = in.hashFileName
	ufcfg.remotepath = path.Base(in.hashFileName)
	wg.Add(1)
	go uploadFile(ufcfg)

	// send the signature file
	ufcfg.localpath = in.hashFileName + sigFileSuffix
	ufcfg.remotepath = path.Base(in.hashFileName) + sigFileSuffix
	wg.Add(1)
	go uploadFile(ufcfg)

	// send each provider zip file
	for k, v := range in.hashesToFilenames {
		ufcfg.localpath = v
		ufcfg.remotepath = path.Base(v)
		wg.Add(1)
		ufcfg.hash = k
		go uploadFile(ufcfg)
	}

	//err = uploadSharedObjects(ctx, in.hashFileName)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = uploadProviderFiles(ctx, hashesToFilenames)
	//if err != nil {
	//	log.Fatal(err)
	//}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs error

	for err = range errChan {
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

type uploadFileIn struct {
	ctx        context.Context
	mgr        *manager.Uploader
	errChan    chan<- error
	localpath  string
	remotepath string
	bucket     string
	hash       string
	wg         *sync.WaitGroup
}

func uploadFile(in uploadFileIn) {
	f, err := os.Open(in.localpath)
	if err != nil {
		in.errChan <- fmt.Errorf("failed to open %q for S3 upload - %w", in.localpath, err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	var checksumSHA256 *string
	if in.hash != "" {
		s1, err := hex.DecodeString(in.hash)
		if err != nil {
			in.errChan <- fmt.Errorf("failed decoding hex string %q - %w", in.hash, err)
			return
		}

		s2 := base64.StdEncoding.EncodeToString(s1)
		checksumSHA256 = &s2
	}

	_, err = in.mgr.Upload(in.ctx, &s3.PutObjectInput{
		Bucket:         &in.bucket,
		Key:            aws.String(path.Base(in.localpath)),
		Body:           f,
		ChecksumSHA256: checksumSHA256,
	})
	if err != nil {
		in.errChan <- fmt.Errorf("failed while uploading %q to s3://%s - %w", in.localpath, in.bucket, err)
		return
	}

	in.errChan <- nil
	in.wg.Done()
}
