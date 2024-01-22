package models

type (
	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		Result string `json:"result"`
	}

	BadResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Shorten struct {
		UUID        int    `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
)
