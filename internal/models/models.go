package models

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}

	ShortenResponse struct {
		Result string `json:"result"`
	}

	BadResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Shorten struct {
		UUID        int64  `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
		UserID      string `json:"user_id"`
	}

	BatchShortenRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	BatchShortenResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	User struct {
		UUID string `json:"uuid"`
	}

	UserURL struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}
)
