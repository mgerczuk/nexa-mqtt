package growatt_web

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type httpClient struct {
	client *http.Client
}

type HttpClient interface {
	postForm(url string, data url.Values, responseBody any) error
}

func (h *httpClient) postForm(url string, data url.Values, responseBody any) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("http.NewRequest failed (web)", slog.String("error", err.Error()))
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	resp, err := h.client.Do(req)
	if err != nil {
		slog.Error("http.Client.Do failed (web)", slog.String("error", err.Error()))
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("io.ReadAll failed (web)", slog.String("error", err.Error()))
		return err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("request failed: (HTTP %s) %s", resp.Status, string(b))
		slog.Error("StatusCode != 200 (web)", slog.String("error", err.Error()))
		return err
	}

	if responseBody != nil {
		if err := json.Unmarshal(b, &responseBody); err != nil {
			slog.Error("json.Unmarshal failed (web)", slog.String("error", err.Error()), slog.String("body", string(b)), slog.String("url", url))
			return err
		}
	}

	return nil
}
