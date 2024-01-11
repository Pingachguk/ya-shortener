package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pingachguk/ya-shortener/internal/models"
)

var storage *FileStorage

func GetStorage() *FileStorage {
	return storage
}

type FileStorage struct {
	f        *os.File
	shortens []models.Shorten
}

func NewFileStorage(filename string) *FileStorage {
	var f *os.File

	storage = &FileStorage{
		shortens: make([]models.Shorten, 0),
	}

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(filename), fs.ModeAppend)
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

		scanner := bufio.NewScanner(f)
		shortens := make([]models.Shorten, 0)
		for scanner.Scan() {
			shorten := models.Shorten{}
			b := scanner.Bytes()
			err := json.Unmarshal(b, &shorten)
			if err != nil {
				panic(err)
			}
			shortens = append(shortens, shorten)
		}
		storage.shortens = shortens
	}

	storage.f = f
	return storage
}

func (fs *FileStorage) AddShorten(shorten models.Shorten) error {
	newID := len(fs.shortens) + 1
	shorten.Uuid = newID

	data, err := shorten.GetJson()
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = fs.f.Write(data)
	if err != nil {
		return err
	}

	fs.shortens = append(fs.shortens, shorten)

	return nil
}
