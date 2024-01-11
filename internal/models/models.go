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
		Uuid        int    `json:"uuid"`
		ShortUrl    string `json:"short_url"`
		OriginalUrl string `json:"original_url"`
	}
)
