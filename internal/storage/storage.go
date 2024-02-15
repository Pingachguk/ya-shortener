package storage

import (
	"context"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/pkg/errors"
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
	AddBatchShorten(ctx context.Context, shortens []models.Shorten) error
	GetByShort(ctx context.Context, short string) (*models.Shorten, error)
	GetByURL(ctx context.Context, URL string) (*models.Shorten, error)
	GetUserURLS(ctx context.Context, userID string) ([]*models.Shorten, error)
	Close(ctx context.Context) error
}

var ErrUnique = errors.New("unique conflict")
