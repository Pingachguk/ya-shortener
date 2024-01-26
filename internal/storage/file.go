package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/pingachguk/ya-shortener/internal/models"
	"os"
	"path/filepath"
)

type FileStorage struct {
	f          *os.File
	countLines int
}

var fileStorage *FileStorage

func InitFileStorage(ctx context.Context, filename string) {
	var f *os.File

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(filename), os.ModeAppend)
		if err != nil {
			panic(err)
		}

		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)

		if err != nil {
			panic(err)
		}
	}

	fileStorage.f = f
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
	scanner := bufio.NewScanner(fs.f)
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
	// ----
	return fs.f.Close()
}

func getCountLines(f *os.File) int {
	numberOfLines := 0
	input := bufio.NewScanner(f)
	for input.Scan() {
		numberOfLines++
	}
	return numberOfLines
}
