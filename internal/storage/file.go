package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type FileStorage struct {
	f          *os.File
	countLines int64
}

var fileStorage *FileStorage

func InitFileStorage(ctx context.Context, path string) {
	var f *os.File

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, os.ModeAppend)
		if err != nil {
			log.Panic().Err(err).Msgf("")
		}

		if dir == path {
			f, err = os.CreateTemp(path, "data_*.json")

			if err != nil {
				log.Panic().Err(err).Msgf("")
			}
		} else {
			f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
			if err != nil {
				log.Panic().Err(err).Msgf("")
			}
		}
	} else if err == nil {
		switch mode := stat.Mode(); {
		case mode.IsDir():
			f, err = os.CreateTemp(path, "data_*.json")

			if err != nil {
				log.Panic().Err(err).Msgf("")
			}
		case mode.IsRegular():
			f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)

			if err != nil {
				log.Panic().Err(err).Msgf("")
			}
		}
	} else {
		log.Panic().Err(err).Msgf("")
	}

	fileStorage = &FileStorage{
		f:          f,
		countLines: getCountLines(f),
	}
}

func GetFileStorage() *FileStorage {
	return fileStorage
}

func (fs *FileStorage) AddShorten(ctx context.Context, shorten models.Shorten) error {
	newID := fs.countLines + 1
	shorten.UUID = newID

	data, err := shorten.GetJSON()
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = fs.f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) GetByShort(ctx context.Context, short string) (*models.Shorten, error) {
	f, err := os.Open(fs.f.Name())
	if err != nil {
		log.Panic().Err(err).Msgf("")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		shorten := models.Shorten{}
		b := scanner.Bytes()
		err := json.Unmarshal(b, &shorten)
		if err != nil {
			panic(err)
		}

		if shorten.ShortURL == short {
			return &shorten, nil
		}
	}

	return nil, nil
}

func (fs *FileStorage) Close(ctx context.Context) error {
	return fs.f.Close()
}

func getCountLines(f *os.File) int64 {
	numberOfLines := 0
	input := bufio.NewScanner(f)
	for input.Scan() {
		numberOfLines++
	}
	return int64(numberOfLines)
}
