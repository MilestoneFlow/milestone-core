package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func initS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"), // Replace with your AWS region
	)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}
	return s3.NewFromConfig(cfg)
}

func uploadFileToS3(s3Client *s3.Client, bucketName, key string, fileData []byte) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   bytes.NewReader(fileData),
	})
	return err
}

func uploadHandler(file *multipart.File, fileHeader *multipart.FileHeader) error {
	s3Client := initS3Client()
	bucketName := "your-bucket-name"
	fileKey := "media_" + uuid.New().String() + filepath.Ext(fileHeader.Filename)

	fileData, err := io.ReadAll(*file)
	if err != nil {
		log.Fatal(err)
	}

	err = uploadFileToS3(s3Client, bucketName, fileKey, fileData)

	return err
}
