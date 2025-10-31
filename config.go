package main

import (
	"os"
	"strconv"
)

var environment *Config

type Config struct {
	MinioEndpoint    string
	MinioAccessKeyID string
	MinioSecretKey   string
	MinioBucket      string
	UploadSecret     []byte
	Email            string
	Password         string
	EmailHost        string
	EmailPort        int
	SelfHost         *string
}

func init(){
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	uploadStorageSecret := os.Getenv("JWT_STORAGE_UPLOAD_SECRET")
	email := os.Getenv("EMAIL_AUTH_USER")
	password := os.Getenv("EMAIL_AUTH_PASS")
	emailHost := os.Getenv("EMAIL_HOST")
	emailPort := os.Getenv("EMAIL_PORT")
	port, _ := strconv.Atoi(emailPort)

	var selfHostPtr *string
	selfHost := os.Getenv("SELF_HOST")
	if selfHost == "" {
		selfHostPtr = nil
	} else {
		selfHostPtr = &selfHost
	}


	environment = &Config{
		MinioEndpoint:    endpoint,
		MinioAccessKeyID: accessKey,
		MinioSecretKey:   secretKey,
		MinioBucket:      bucket,
		UploadSecret:     []byte(uploadStorageSecret),
		Email:            email,
		Password:         password,
		EmailHost:        emailHost,
		EmailPort:        port,
		SelfHost:         selfHostPtr,
	}
}

func GetConfig() Config {
	return *environment
}
