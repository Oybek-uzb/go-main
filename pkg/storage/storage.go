package storage

import (
	"context"
)

type UploadInput struct {
	File        string
	Name        string
	Folder 		string
}

type Storage interface {
	Upload(ctx context.Context, input UploadInput) (string, error)
}
