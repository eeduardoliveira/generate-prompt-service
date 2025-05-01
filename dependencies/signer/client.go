package signer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client interface {
	GetPresignedUploadURL(ctx context.Context, bucket, clienteID string) (string, error)
}

type HTTPClient struct {
	BaseURL string
	HTTP    *http.Client
}

type signedURLResponse struct {
	URL string `json:"url"`
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		BaseURL: os.Getenv("BUCKET_SIGNER_URL"),
		HTTP:    &http.Client{},
	}
}

func (c *HTTPClient) GetPresignedUploadURL(ctx context.Context, bucket, clienteID string) (string, error) {
	endpoint := fmt.Sprintf("%s/signed-url?bucket=%s&clienteID=%s&upload=true", c.BaseURL, bucket, clienteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao chamar bucket-signer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bucket-signer retornou status %d", resp.StatusCode)
	}

	var result signedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta do bucket-signer: %w", err)
	}

	return result.URL, nil
}