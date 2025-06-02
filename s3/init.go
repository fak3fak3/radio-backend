package s3

import (
	"context"
	"go-postgres-gorm-gin-api/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Client *minio.Client
}

var S3Instance *S3Client

func ConnectMinio(cfg *config.Config) (*S3Client, error) {
	var err error
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	S3Instance.Client = client

	log.Println("Connected to MinIO server")

	// Check if the bucket exists
	bucketExists, err := S3Instance.Client.BucketExists(context.Background(), cfg.MinioBucketName)
	if err != nil {
		log.Fatalln(err)
	}
	if !bucketExists {
		// Create the bucket if it doesn't exist
		err = S3Instance.Client.MakeBucket(context.Background(), cfg.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Bucket %s created successfully\n", cfg.MinioBucketName)
	} else {
		log.Printf("Bucket %s already exists\n", cfg.MinioBucketName)
	}

	return S3Instance, nil
}
