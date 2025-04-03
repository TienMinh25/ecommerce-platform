package pkg

import (
	"context"
	"net/http"
)

// RequestOption là hàm tùy chọn để cấu hình request
type RequestOption func(*http.Request) error

type ResponseAPI struct {
	StatusCode int
	Headers    http.Header
	RawBody    []byte
}

type HTTPClient interface {
	// SendRequest Sends an HTTP request using the specified method (e.g., GET, POST), URL, and optional data payload (`io.Reader`).
	// It returns the response body as a byte array or an error if the request fails.
	SendRequest(ctx context.Context, method string, url string, options ...RequestOption) (*ResponseAPI, error)
}
