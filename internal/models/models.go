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
)
