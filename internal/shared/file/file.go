package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"
	"sync"

	"github.com/maximfedotov74/cloud-api/internal/shared/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var once sync.Once
var minioClient *minio.Client

type FileClient struct {
	minio *minio.Client
}

type UploadResponse struct {
	Path string `json:"path"`
}

// Получение размера bucket
//bucketInfo, err := minioClient.BucketInfo(ctx, bucketName)

//make link
//presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, nil

// GetObject

func (c *FileClient) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := c.minio.BucketExists(ctx, bucketName)

	if err != nil {
		return fmt.Errorf("failed when find bucket, cause: %s", err.Error())
	}

	if exists {
		return fmt.Errorf("bucket alreay exists")
	}

	err = c.minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
	if err != nil {
		return fmt.Errorf("failed when make new bucket, cause: %s", err.Error())
	}

	return nil
}

func New(minioUrl string, user string, password string) *FileClient {
	once.Do(func() {
		client, err := minio.New(minioUrl, &minio.Options{Creds: credentials.NewStaticV4(user, password, ""), Secure: false})
		if err != nil {
			log.Fatalf("Failed to connect to minio service, cause: %s", err.Error())
		}
		minioClient = client
	})
	log.Println("File manager Minio successfully connected!")
	return &FileClient{minio: minioClient}
}

func (c *FileClient) Upload(ctx context.Context, bucketName string, fileName string, h *multipart.FileHeader) (*UploadResponse, error) {
	file, err := h.Open()
	if err != nil {
		return nil, fmt.Errorf("error when open file, cause: %s", err.Error())
	}
	defer file.Close()
	contentType := h.Header.Get("Content-type")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error when get file bytes with io: %s", err.Error())
	}
	// splittedContentType := strings.Split(contentType, "/")
	// fileType := splittedContentType[0]
	// extType := splittedContentType[1]
	ext := utils.GetFileExt(h.Filename)
	newName := fileName + ext
	reader := bytes.NewReader(fileBytes)
	//change empty string to bucket name
	_, err = c.minio.PutObject(ctx, bucketName, newName, reader, reader.Size(), minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	})
	if err != nil {
		return nil, fmt.Errorf("error when uploading file, cause: %s", err.Error())
	}
	return &UploadResponse{Path: path.Join("/", "storage", bucketName, newName)}, nil
}
