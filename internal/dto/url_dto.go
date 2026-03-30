package dto

// CreateURLRequest is the request body for creating a short URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url"`
	CustomCode  string `json:"custom_code,omitempty"`
}

// CreateURLResponse is the response body after creating a short URL
type CreateURLResponse struct {
	ShortCode      string `json:"short_code"`
	ShortURL       string `json:"short_url"`
	OriginalURL    string `json:"original_url"`
	URLAvailable   bool   `json:"url_available"`
	SEOTitle       string `json:"seo_title,omitempty"`
	SEODescription string `json:"seo_description,omitempty"`
	SEOImage       string `json:"seo_image,omitempty"`
	Warning        string `json:"warning,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
