package storage

import (
	"abir/pkg/utils"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
	"image/jpeg"
	"image/png"
	"strings"
	"time"
)

const (
	timeout = time.Second * 20
)

type FileStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
	env      string
}

func NewFileStorage(client *minio.Client, bucket, endpoint, env string) *FileStorage {
	return &FileStorage{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
		env:      env,
	}
}

func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	data := input.File
	i := strings.Index(data, ",")
	j := strings.Index(data, ";")
	k := strings.Index(data, ":")
	if i < 0 || j < 0 || k < 0 {
		return "", errors.New("invalid file format")
	}
	contentType := data[k+1:j]
	l := strings.Index(contentType, "/")
	if l < 0 {
		return "", errors.New("invalid file format")
	}
	fileType := contentType[l+1:]
	if fileType != "jpeg" && fileType != "png" {
		return "", errors.New("file must be png or jpeg")
	}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[i+1:]))
	opts := minio.PutObjectOptions{
		ContentType:  contentType,
	}

	newCtx, clFn := context.WithTimeout(ctx, timeout)
	defer clFn()

	if fileType == "png" {
		imgPng, err := png.Decode(dec)
		if err != nil {
			return "", err
		}
		originalImage := resize.Resize(1024, 0, imgPng, resize.Lanczos3)
		originalImageBase64, err := utils.ImageToReader(originalImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "original", input.Name, fileType), originalImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}
		mediumImage := resize.Resize(512, 0, imgPng, resize.Lanczos3)
		mediumImageBase64, err := utils.ImageToReader(mediumImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "medium", input.Name, fileType), mediumImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}
		smallImage := resize.Resize(256, 0, imgPng, resize.Lanczos3)
		smallImageBase64, err := utils.ImageToReader(smallImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "small", input.Name, fileType), smallImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}
	}
	if fileType == "jpeg" {
		imgJpg, err := jpeg.Decode(dec)
		if err != nil {
			return "", err
		}
		originalImage := resize.Resize(1024, 0, imgJpg, resize.Lanczos3)
		originalImageBase64, err := utils.ImageToReader(originalImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "original", input.Name, fileType), originalImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}

		mediumImage := resize.Resize(512, 0, imgJpg, resize.Lanczos3)
		mediumImageBase64, err := utils.ImageToReader(mediumImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "medium", input.Name, fileType), mediumImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}
		smallImage := resize.Resize(256, 0, imgJpg, resize.Lanczos3)
		smallImageBase64, err := utils.ImageToReader(smallImage, fileType)
		if err != nil {
			return "", err
		}
		_, err = fs.client.PutObjectWithContext(newCtx,
			fs.bucket, fmt.Sprintf("%s/%s/%s.%s", input.Folder, "small", input.Name, fileType), smallImageBase64, -1, opts)
		if err != nil {
			logrus.Errorf("error occured while uploading file to bucket: %s", err.Error())
			return "", err
		}
	}

	return fmt.Sprintf("%s/%s.%s", input.Folder, input.Name, fileType), nil
}
