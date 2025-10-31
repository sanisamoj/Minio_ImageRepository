package main

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinIO(config *Config) (*minio.Client, error) {
	ctx := context.Background()
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKeyID, config.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func EnsureBucketExists(client *minio.Client, bucketName string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
