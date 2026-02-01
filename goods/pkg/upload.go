package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"log"
	"mime/multipart"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func Upload(file *multipart.FileHeader) string {
	objectName := file.Filename
	region := "us-east-1"
	access_key_id := "Ctxwr2HEaMLQjo0Buh5z"
	secret_access_key := "fHlLmXCFBJ82iSrjxea0KndQvIWRsq4tYMb5wV76"
	endpoint := "http://115.190.54.31:9501"

	// build aws.Config
	cfg := aws.Config{
		Region: region,
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: endpoint,
			}, nil
		}),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access_key_id, secret_access_key, "")),
	}

	// build S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	ext := strings.Split(objectName, ".")[1]
	objectName = fmt.Sprintf("%d.%s", time.Now().UnixMilli(), ext)
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("rui"),
		Key:    aws.String(objectName),
		Body:   strings.NewReader("hello rustfs"),
	})
	if err != nil {
		log.Fatalf("上传文件失败: %v", err)
	}
	return fmt.Sprintf("文件上传成功：%s", objectName)
}
