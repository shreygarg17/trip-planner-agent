package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestSearchWeather_Success(t *testing.T) {
	// Reset weather singleton for the test
	weatherOnce = sync.Once{}
	weatherInstance = nil

	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.Host, "geocoding-api") {
					resp := httptest.NewRecorder()
					resp.Header().Set("Content-Type", "application/json")
					resp.WriteString(`{"results": [{"latitude": 27.5, "longitude": 88.5}]}`)
					return resp.Result(), nil
				}
				if strings.Contains(req.URL.Host, "api.open-meteo") {
					resp := httptest.NewRecorder()
					resp.Header().Set("Content-Type", "application/json")
					resp.WriteString(`{
						"daily": {
							"temperature_2m_max": [25.0, 26.0, 24.5, 25.5, 26.5, 27.0, 25.0],
							"temperature_2m_min": [15.0, 16.0, 14.5, 15.5, 16.5, 17.0, 15.0],
							"precipitation_probability_max": [20, 30, 40, 10, 20, 15, 25]
						}
					}`)
					return resp.Result(), nil
				}
				return nil, nil
			},
		},
		Timeout: 2 * time.Second,
	}

	ws := NewOpenMeteoClient(mockClient)
	info, err := ws.SearchWeather(context.Background(), "Sikkim")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if info.Condition != "Moderate" {
		t.Errorf("Expected condition 'Moderate', got: %s", info.Condition)
	}

	if info.TemperatureMax < 25.0 || info.TemperatureMax > 26.0 {
		t.Errorf("Unexpected average max temperature: %f", info.TemperatureMax)
	}
}

func TestSearchWeather_Integration(t *testing.T) {
	// Reset weather singleton for the test
	weatherOnce = sync.Once{}
	weatherInstance = nil

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws := NewOpenMeteoClient(nil)
	info, err := ws.SearchWeather(ctx, "Gangtok")
	if err != nil {
		t.Logf("[SKIPPED] Integration test call to Open-Meteo failed (offline or slow): %v", err)
		return
	}

	if info.Condition == "" {
		t.Error("Expected condition to be populated in integration test")
	}
}
