package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	exif "github.com/usagiga/go-exif-remove"
	"log"
)

func main() {
	lambda.Start(OnObjectCreated)
}

func OnObjectCreated(ctx context.Context, ev events.S3Event) (err error) {
	bucket := ev.Records[0].S3.Bucket.Name
	objKey := ev.Records[0].S3.Object.Key
	sess := session.Must(session.NewSession()) // credential through IAM Role

	file, acl, err := Download(sess, bucket, objKey)
	if err != nil {
		return fmt.Errorf("failed download: %w", err)
	}

	file, err = RemoveExif(file)
	if err != nil {
		// Not to run redundantly, guard it
		if errors.Is(err, exif.ErrNoExif) ||
			errors.Is(err, exif.ErrNotCompatible) {
			return
		}

		return fmt.Errorf("failed remove exif: %w", err)
	}


	err = Upload(sess, bucket, objKey, acl, file)
	if err != nil {
		return fmt.Errorf("failed upload: %w", err)
	}

	log.Printf("removed %s exif", objKey)
	return nil
}

func Download(cp client.ConfigProvider, bucket, objKey string) (fileBytes []byte, acl string, err error) {
	// Get file
	buf := aws.NewWriteAtBuffer(nil)
	downloader := s3manager.NewDownloader(cp)
	_, err = downloader.Download(buf, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objKey,
	})
	if err != nil {
		return nil, "", err
	}

	// Get ACL
	s3svc := s3.New(cp)
	getAclOut, err := s3svc.GetObjectAcl(&s3.GetObjectAclInput{
		Bucket:              &bucket,
		Key:                 &objKey,
	})
	if err != nil {
		return nil, "", err
	}
	acl = getAclOut.String()

	return buf.Bytes(), acl, nil
}

func Upload(cp client.ConfigProvider, bucket, objKey, acl string, fileBytes []byte) (err error) {
	buf := bytes.NewBuffer(fileBytes)
	uploader := s3manager.NewUploader(cp)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   buf,
		Bucket: &bucket,
		Key:    &objKey,
		ACL:    &acl,
	})
	if err != nil {
		return err
	}

	return nil
}

func RemoveExif(imgBytes []byte) (removedBytes []byte, err error) {
	removedBytes, err = exif.Remove(imgBytes)
	if err != nil {
		return nil, err
	}

	return removedBytes, nil
}
