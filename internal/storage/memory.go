package storage

import (
	"context"
	"errors"
	"github.com/pingachguk/ya-shortener/internal/models"
)

type MemoryStorage struct {
	shortens []models.Shorten
	close    bool
}

var memoryStorage *MemoryStorage

func InitMemoryStorage() {
	memoryStorage = &MemoryStorage{
		shortens: make([]models.Shorten, 0),
	}
}

func GetMemoryStorage() *MemoryStorage {
	return memoryStorage
}

func (ms *MemoryStorage) AddShorten(ctx context.Context, shorten models.Shorten) error {
	newID := len(ms.shortens) + 1
	shorten.UUID = int64(newID)

	ms.shortens = append(ms.shortens, shorten)

	return nil
}

func (ms *MemoryStorage) AddBatchShorten(ctx context.Context, shortens []models.Shorten) error {
	for _, shorten := range shortens {
		ms.AddShorten(ctx, shorten)
	}

	return nil
}

func (ms *MemoryStorage) GetByShort(ctx context.Context, short string) (*models.Shorten, error) {
	if ms.close {
		return nil, errors.New("storage is closed")
	}

	for _, v := range ms.shortens {
		if v.ShortURL == short {
			return &v, nil
		}
	}

	return nil, nil
}

func (ms *MemoryStorage) Close(ctx context.Context) error {
	ms.close = true

	return nil
}
