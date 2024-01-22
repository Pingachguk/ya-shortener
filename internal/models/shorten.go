package models

import "encoding/json"

func NewShorten(shortURL string, originalURL string) *Shorten {
	return &Shorten{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func (s Shorten) GetJSON() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}
