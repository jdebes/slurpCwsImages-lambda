package service

import (
	"time"
	"context"
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

type AwsService interface {
	UploadItemToS3(file []byte, format string, key string) error
}

type awsServiceImpl struct {
	bucket string
	timeout time.Duration
	svc *s3.S3
}

func (service *awsServiceImpl) UploadItemToS3(file []byte, format string, key string) error {
	var cancelFn func()
	ctx := context.Background()
	fileStream := bytes.NewReader(file)

	if service.timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, service.timeout)
	}
	defer cancelFn()

	_, err := service.svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(service.bucket),
		Key:    aws.String(key),
		Body:   fileStream,
		ContentLength: aws.Int64(fileStream.Size()),
		ContentType: aws.String(fmt.Sprintf("image/%s", format)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			log.WithError(err).Error("Failed due to a cancelled context")
		} else {
			log.WithError(err).Error("Unable to upload to S3 bucket")
		}

		return err
	}

	return nil
}

func BuildAwsService(bucket string, timeout time.Duration) AwsService {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.ApSoutheast2RegionID),


	})

	return &awsServiceImpl{bucket: bucket, timeout: timeout, svc: svc}
}