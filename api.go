package akamai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type SensorInput struct {
	Abck      string `json:"_abck"`
	Bmsz      string `json:"bm_sz"`
	PageUrl   string `json:"pageUrl"`
	UserAgent string `json:"ua"`
	ApiKey    string `json:"apiKey"`
}

// GenerateSensorData returns the sensor data required to generate valid akamai cookies using the FRIS Solutions API.
func (s *Session) GenerateSensorData(ctx context.Context, input *SensorInput) (string, error) {
	const sensorEndpoint = "https://api.frisapi.dev/akamai/sensor"
	input.ApiKey = s.apiKey
	return s.sendRequest(ctx, sensorEndpoint, input)
}

type PixelInput struct {
	UserAgent string `json:"userAgent"`
	PageUrl   string `json:"pageUrl"`
	PixelId   string `json:"pixelId"`
	ScriptVar string `json:"scriptVar"`
	ApiKey    string `json:"apiKey"`
}

// GeneratePixelData returns the pixel data using the FRIS Solutions API.
func (s *Session) GeneratePixelData(ctx context.Context, input *PixelInput) (string, error) {
	const pixelEndpoint = "https://api.frisapi.dev/akamai/pixel"
	input.ApiKey = s.apiKey
	return s.sendRequest(ctx, pixelEndpoint, input)
}

func (s *Session) sendRequest(ctx context.Context, url string, input any) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("missing api key")
	}
	payload, err := jsoniter.Marshal(input)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept-encoding", "gzip")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Data    string            `json:"data,omitempty"`
		Error   string            `json:"error,omitempty"`
		Headers map[string]string `json:"headers,omitempty"`
		Success bool              `json:"success"`
	}
	if err := jsoniter.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	if response.Error != "" {
		return "", fmt.Errorf("api returned with: %s", response.Error)
	}

	return response.Data, nil
}
