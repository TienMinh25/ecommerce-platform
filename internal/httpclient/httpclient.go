package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type httpClient struct {
	tracer pkg.Tracer
	client *http.Client
}

func NewHTTPClient(tracer pkg.Tracer) pkg.HTTPClient {
	return &httpClient{
		tracer: tracer,
		client: &http.Client{
			Timeout: time.Second * 40,
		},
	}
}

func (h *httpClient) SendRequest(ctx context.Context, method string, url string, options ...pkg.RequestOption) (*pkg.ResponseAPI, error) {
	ctx, span := h.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.AdapterLayer, "SendRequest"))
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for _, option := range options {
		if err = option(req); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	resp, err := h.client.Do(req)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &pkg.ResponseAPI{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		RawBody:    body,
	}, nil
}

// WithHeader thêm header vào request
func WithHeader(key, value string) pkg.RequestOption {
	return func(req *http.Request) error {
		req.Header.Add(key, value)
		return nil
	}
}

// WithHeaders thêm nhiều headers vào request
func WithHeaders(headers map[string]string) pkg.RequestOption {
	return func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		return nil
	}
}

// WithJSONBody thêm JSON body vào request
func WithJSONBody[T any](body T) pkg.RequestOption {
	return func(req *http.Request) error {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		return nil
	}
}

// WithBearerToken thêm bearer token vào request
func WithBearerToken(token string) pkg.RequestOption {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// WithFormBody thêm form data vào request
func WithFormBody(form map[string]string) pkg.RequestOption {
	return func(req *http.Request) error {
		values := url.Values{}
		for key, value := range form {
			values.Add(key, value)
		}

		// Encode data form
		formEncoded := values.Encode()

		// Tạo body cho request
		req.Body = io.NopCloser(strings.NewReader(formEncoded))

		// Đặt Content-Type cho form data
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Đặt Content-Length
		req.ContentLength = int64(len(formEncoded))

		return nil
	}
}
