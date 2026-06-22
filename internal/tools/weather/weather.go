package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/shreygarg/trip-planner-agent/constants"
	"github.com/shreygarg/trip-planner-agent/models"
	"github.com/shreygarg/trip-planner-agent/utils"
)

// WeatherService defines the contract for fetching weather information.
type WeatherService interface {
	SearchWeather(ctx context.Context, destination string) (*models.WeatherInfo, error)
}

var (
	weatherOnce     sync.Once
	weatherInstance *OpenMeteoClient
)

// OpenMeteoClient implements WeatherService calling Open-Meteo APIs.
type OpenMeteoClient struct {
	client *http.Client
}

// NewOpenMeteoClient instantiates a new OpenMeteoClient (singleton).
func NewOpenMeteoClient(client *http.Client) *OpenMeteoClient {
	weatherOnce.Do(func() {
		if client == nil {
			client = &http.Client{
				Timeout: constants.HTTPTimeout,
			}
		}
		weatherInstance = &OpenMeteoClient{
			client: client,
		}
	})
	return weatherInstance
}

// GetCoordinates retrieves coordinates for a destination name.
func (c *OpenMeteoClient) GetCoordinates(ctx context.Context, destination string) (*models.Coordinates, error) {
	reqID := utils.GetRequestID(ctx)
	httpCtx, cancel := context.WithTimeout(ctx, constants.HTTPTimeout)
	defer cancel()

	apiURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s", url.QueryEscape(destination))
	req, err := http.NewRequestWithContext(httpCtx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	log.Printf("[%s] [INFO] Geocoding location: %s", reqID, destination)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocoding api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding api returned unexpected status: %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(data.Results) == 0 {
		return nil, fmt.Errorf("no coordinates found for: %s", destination)
	}

	return &models.Coordinates{
		Latitude:  data.Results[0].Latitude,
		Longitude: data.Results[0].Longitude,
	}, nil
}

// GetForecast retrieves weather info using coordinates.
func (c *OpenMeteoClient) GetForecast(ctx context.Context, latitude, longitude float64) (*models.WeatherInfo, error) {
	reqID := utils.GetRequestID(ctx)
	httpCtx, cancel := context.WithTimeout(ctx, constants.HTTPTimeout)
	defer cancel()

	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_max&forecast_days=7",
		latitude,
		longitude,
	)

	req, err := http.NewRequestWithContext(httpCtx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast request: %w", err)
	}

	log.Printf("[%s] [INFO] Fetching forecast for coordinates: %f, %f", reqID, latitude, longitude)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forecast api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast api returned unexpected status: %d", resp.StatusCode)
	}

	var data struct {
		Daily struct {
			Temperature2mMax            []float64 `json:"temperature_2m_max"`
			Temperature2mMin            []float64 `json:"temperature_2m_min"`
			PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
		} `json:"daily"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse forecast response: %w", err)
	}

	return c.mapForecastToWeatherInfo(&data)
}

// SearchWeather resolves coordinates and gets the forecast.
func (c *OpenMeteoClient) SearchWeather(ctx context.Context, destination string) (*models.WeatherInfo, error) {
	coords, err := c.GetCoordinates(ctx, destination)
	if err != nil {
		return nil, err
	}
	return c.GetForecast(ctx, coords.Latitude, coords.Longitude)
}

func (c *OpenMeteoClient) mapForecastToWeatherInfo(data *struct {
	Daily struct {
		Temperature2mMax            []float64 `json:"temperature_2m_max"`
		Temperature2mMin            []float64 `json:"temperature_2m_min"`
		PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
	} `json:"daily"`
}) (*models.WeatherInfo, error) {
	daily := data.Daily
	n := len(daily.Temperature2mMax)
	if n == 0 {
		return nil, fmt.Errorf("empty forecast daily data")
	}

	var sumMax, sumMin, sumRain float64
	for i := 0; i < n; i++ {
		sumMax += daily.Temperature2mMax[i]
		sumMin += daily.Temperature2mMin[i]
		sumRain += daily.PrecipitationProbabilityMax[i]
	}

	avgMax := sumMax / float64(n)
	avgMin := sumMin / float64(n)
	avgRain := sumRain / float64(n)

	condition := "Moderate"
	if avgRain > 50.0 {
		condition = "Rainy"
	} else if avgMax > 30.0 {
		condition = "Sunny"
	} else if avgMax < 10.0 {
		condition = "Cold"
	}

	return &models.WeatherInfo{
		TemperatureMax:  avgMax,
		TemperatureMin:  avgMin,
		RainProbability: avgRain,
		Condition:       condition,
	}, nil
}
