package models

import "encoding/json"

func NewShorten(shortURL string, originalURL string, userID string) *Shorten {
	return &Shorten{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
}

func (s Shorten) GetJSON() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}
