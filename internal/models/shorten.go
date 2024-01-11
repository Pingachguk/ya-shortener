package models

import "encoding/json"

func NewShorten(shortUrl string, originalUrl string) *Shorten {
	return &Shorten{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
	}
}

func (s Shorten) GetJSON() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}
