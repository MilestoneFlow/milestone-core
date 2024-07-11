package awsinternal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"mime"
	"path/filepath"
)

func GetConfiguration(region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider("AKIAZI2LEQP6AVAFF6XR", "PKjulG2Kz6ua+rW3gE17sOMq0xJehrdyC30bUhAw", ""),
	), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func UploadToS3(ctx context.Context, filename string, fileData io.Reader) error {
	cfg, err := GetConfiguration("us-east-1")
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(*cfg)

	bucket := "milestone-uploaded-flows-media"
	key := "step_media/" + filename

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        fileData,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(filename))),
	})
	return err
}
