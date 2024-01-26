package storage

import (
	"context"
	"github.com/pingachguk/ya-shortener/internal/models"
)

func GetStorage() Storage {
	if database != nil {
		return database
	} else if fileStorage != nil {
		return fileStorage
	} else if memoryStorage != nil {
		return memoryStorage
	}

	return nil
}

type Storage interface {
	AddShorten(ctx context.Context, shorten models.Shorten) error
	GetByShort(ctx context.Context, short string) (*models.Shorten, error)
	Close(ctx context.Context) error
}
