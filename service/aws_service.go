package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

const (
	AwsNotFound = "NotFound"
)

type AwsService interface {
	UploadItemToS3(file []byte, fileExt string, productID string, imageType string) error
	S3ItemExists(productID string, fileExt string, imageType string) bool
}

type awsServiceImpl struct {
	bucket  string
	timeout time.Duration
	svc     *s3.S3
}

func (service *awsServiceImpl) UploadItemToS3(file []byte, fileExt string, productID string, imageType string) error {
	fileStream := bytes.NewReader(file)

	ctx, cancelFn := context.WithTimeout(context.Background(), service.timeout)
	defer cancelFn()

	startTime := time.Now()
	_, err := service.svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(service.bucket),
		Key:           aws.String(service.buildProductFileName(productID, fileExt, imageType)),
		Body:          fileStream,
		ContentLength: aws.Int64(fileStream.Size()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			log.WithError(err).Error("Failed due to a cancelled context")
		} else {
			log.WithError(err).Error("Unable to upload to S3 bucket")
		}

		return err
	}
	duration := time.Now().Sub(startTime)

	log.Debug(fmt.Sprintf("Uploaded to S3 %s", duration.String()))
	return nil
}

func (service *awsServiceImpl) S3ItemExists(productID string, fileExt string, imageType string) bool {
	ctx, cancelFn := context.WithTimeout(context.Background(), service.timeout)
	defer cancelFn()

	_, err := service.svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(service.bucket),
		Key:    aws.String(service.buildProductFileName(productID, fileExt, imageType)),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)

		if ok && aerr.Code() == AwsNotFound {
			return false
		}

		if ok && aerr.Code() == request.CanceledErrorCode {
			log.WithError(err).Error("Failed due to a cancelled context")
		} else {
			log.WithError(err).Error("Unable to upload to S3 bucket")
		}
	}

	return true
}

func (service *awsServiceImpl) buildProductFileName(productID string, fileExt string, imageType string) string {
	return fmt.Sprintf("%s_%s%s", productID, strings.ToLower(imageType), fileExt)
}

func BuildAwsService(bucket string, timeout time.Duration) AwsService {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region: aws.String(endpoints.ApSoutheast2RegionID),
	})

	return &awsServiceImpl{bucket: bucket, timeout: timeout, svc: svc}
}
