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
			log.Panic().Err(err).Msgf("error make dir")
		}

		if dir == path {
			f, err = os.CreateTemp(path, "data_*.json")

			if err != nil {
				log.Panic().Err(err).Msgf("error create temp")
			}
		} else {
			f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
			if err != nil {
				log.Panic().Err(err).Msgf("error open file")
			}
		}
	} else if err == nil {
		switch mode := stat.Mode(); {
		case mode.IsDir():
			f, err = os.CreateTemp(path, "data_*.json")

			if err != nil {
				log.Panic().Err(err).Msgf("error create temp")
			}
		case mode.IsRegular():
			f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)

			if err != nil {
				log.Panic().Err(err).Msgf("error open file")
			}
		}
	} else {
		log.Panic().Err(err).Msgf("bad get stat by path")
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

func (fs *FileStorage) AddBatchShorten(ctx context.Context, shortens []models.Shorten) error {
	data := make([]byte, 0)
	for _, shorten := range shortens {
		newID := fs.countLines + 1
		fs.countLines = newID
		shorten.UUID = newID

		row, err := shorten.GetJSON()
		if err != nil {
			return err
		}

		row = append(row, '\n')
		data = append(data, row...)
	}

	_, err := fs.f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) GetByShort(ctx context.Context, short string) (*models.Shorten, error) {
	f, err := os.Open(fs.f.Name())
	if err != nil {
		log.Panic().Err(err).Msgf("error open file for read")
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

func (fs *FileStorage) GetByURL(ctx context.Context, URL string) (*models.Shorten, error) {
	f, err := os.Open(fs.f.Name())
	if err != nil {
		log.Panic().Err(err).Msgf("error opem file")
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

		if shorten.OriginalURL == URL {
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
