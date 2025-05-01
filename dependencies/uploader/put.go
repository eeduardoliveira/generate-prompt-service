package uploader

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	UploadJSON(ctx context.Context, presignedURL string, body io.Reader) error
}

type HTTPUploader struct {
	Client *http.Client
}

func NewHTTPUploader() *HTTPUploader {
	return &HTTPUploader{
		Client: &http.Client{},
	}
}

func (u *HTTPUploader) UploadJSON(ctx context.Context, presignedURL string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, presignedURL, body)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição PUT: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := u.Client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição PUT: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("falha no upload, status %d", resp.StatusCode)
	}

	return nil
}