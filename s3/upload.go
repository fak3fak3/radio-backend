package s3

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

func (c *S3Client) UploadFile(fileHeader *multipart.FileHeader, bucketName, objectName string, fileType string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = c.Client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileType,
		},
	)
	return err
}
