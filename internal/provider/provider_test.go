package provider

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hwang/app-icon-cli/internal/config"
)

func TestHTTPClientGenerateDecodesOpenAIStyleResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer key" {
			t.Fatalf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"b64_json":"` + base64.StdEncoding.EncodeToString([]byte("png")) + `"}]}`))
	}))
	defer server.Close()

	client := HTTPClient{Client: server.Client()}
	images, err := client.Generate(context.Background(), GenerateRequest{
		Provider: "test",
		Prompt:   "prompt",
		Count:    1,
		Size:     1024,
		Config:   config.ProviderConfig{APIKey: "key", BaseURL: server.URL, Model: "m"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(images[0].Data) != "png" {
		t.Fatalf("image data = %q", string(images[0].Data))
	}
}
