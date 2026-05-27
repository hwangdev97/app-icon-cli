package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hwang/app-icon-cli/internal/config"
)

type Image struct {
	Data []byte
	URL  string
}

type Client interface {
	Generate(ctx context.Context, req GenerateRequest) ([]Image, error)
}

type GenerateRequest struct {
	Provider string                `json:"provider"`
	Prompt   string                `json:"prompt"`
	Count    int                   `json:"count"`
	Size     int                   `json:"size"`
	Config   config.ProviderConfig `json:"-"`
}

type HTTPClient struct {
	Client *http.Client
}

func NewHTTP() HTTPClient {
	return HTTPClient{Client: &http.Client{Timeout: 90 * time.Second}}
}

func (c HTTPClient) Generate(ctx context.Context, req GenerateRequest) ([]Image, error) {
	apiKey, err := config.ResolveAPIKey(req.Config.APIKey)
	if err != nil {
		return nil, err
	}
	if apiKey == "" {
		return nil, fmt.Errorf("provider %s api key is not configured", req.Provider)
	}
	if req.Config.BaseURL == "" {
		return nil, fmt.Errorf("provider %s base_url is not configured", req.Provider)
	}
	body := map[string]any{
		"model":  req.Config.Model,
		"prompt": req.Prompt,
		"n":      req.Count,
		"size":   fmt.Sprintf("%dx%d", req.Size, req.Size),
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, req.Config.BaseURL, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("provider %s returned %s: %s", req.Provider, resp.Status, string(respBody))
	}
	var decoded struct {
		Images []struct {
			B64 string `json:"b64_json"`
			URL string `json:"url"`
		} `json:"images"`
		Data []struct {
			B64 string `json:"b64_json"`
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return nil, fmt.Errorf("decode provider response: %w", err)
	}
	items := decoded.Images
	if len(items) == 0 {
		items = decoded.Data
	}
	images := make([]Image, 0, len(items))
	for _, item := range items {
		img := Image{URL: item.URL}
		if item.B64 != "" {
			data, err := base64.StdEncoding.DecodeString(item.B64)
			if err != nil {
				return nil, fmt.Errorf("decode image data: %w", err)
			}
			img.Data = data
		}
		images = append(images, img)
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("provider %s response did not include images", req.Provider)
	}
	return images, nil
}

func (c HTTPClient) client() *http.Client {
	if c.Client != nil {
		return c.Client
	}
	return http.DefaultClient
}
