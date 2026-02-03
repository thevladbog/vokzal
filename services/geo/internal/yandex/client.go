// Package yandex — клиент Yandex Geocoder API.
package yandex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// Client — клиент для Yandex Geocoder API.
type Client struct {
	client  *http.Client
	logger  *zap.Logger
	apiKey  string
	baseURL string
}

// GeocodeResponse — ответ Yandex Geocoder API.
type GeocodeResponse struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Address struct {
								Formatted string `json:"formatted"`
							} `json:"Address"`
						} `json:"GeocoderMetaData"`
					} `json:"metaDataProperty"`
					Point struct {
						Pos string `json:"pos"` // "lon lat"
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

// GeocodeResult — результат геокодирования.
type GeocodeResult struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NewClient создаёт новый Client.
func NewClient(apiKey, baseURL string, logger *zap.Logger) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Geocode преобразует адрес в координаты.
func (c *Client) Geocode(address string) (*GeocodeResult, error) {
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("geocode", address)
	params.Set("format", "json")
	params.Set("results", "1")

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	c.logger.Debug("Geocoding address", zap.String("address", address))

	resp, err := c.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Debug("failed to close response body", zap.Error(closeErr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result GeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Response.GeoObjectCollection.FeatureMember) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	geoObject := result.Response.GeoObjectCollection.FeatureMember[0].GeoObject
	formatted := geoObject.MetaDataProperty.GeocoderMetaData.Address.Formatted
	pos := geoObject.Point.Pos

	var lon, lat float64
	if _, err := fmt.Sscanf(pos, "%f %f", &lon, &lat); err != nil {
		return nil, fmt.Errorf("failed to parse coordinates: %w", err)
	}

	c.logger.Info("Geocoded address", zap.String("address", formatted), zap.Float64("lat", lat), zap.Float64("lon", lon))

	return &GeocodeResult{
		Address:   formatted,
		Latitude:  lat,
		Longitude: lon,
	}, nil
}

// ReverseGeocode преобразует координаты в адрес.
func (c *Client) ReverseGeocode(lat, lon float64) (*GeocodeResult, error) {
	geocode := fmt.Sprintf("%f,%f", lon, lat)

	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("geocode", geocode)
	params.Set("format", "json")
	params.Set("results", "1")

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	c.logger.Debug("Reverse geocoding", zap.Float64("lat", lat), zap.Float64("lon", lon))

	resp, err := c.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Debug("failed to close response body", zap.Error(closeErr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result GeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Response.GeoObjectCollection.FeatureMember) == 0 {
		return nil, fmt.Errorf("no results found for coordinates: %f, %f", lat, lon)
	}

	formatted := result.Response.GeoObjectCollection.FeatureMember[0].GeoObject.MetaDataProperty.GeocoderMetaData.Address.Formatted

	c.logger.Info("Reverse geocoded", zap.String("address", formatted))

	return &GeocodeResult{
		Address:   formatted,
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
