package dto

// CreateURLRequest is the request body for creating a short URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url"`
}

// CreateURLResponse is the response body after creating a short URL
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
