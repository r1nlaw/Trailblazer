package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"trailblazer/internal/config"
)

type GeocodingClient struct {
	baseURL string
	token   string
	Client  *http.Client
}

func NewGeocodingClient(cfg config.GeocoderConfig) *GeocodingClient {
	return &GeocodingClient{
		baseURL: cfg.Geo_url,
		token:   cfg.Api_key,
		Client:  http.DefaultClient,
	}
}

func (c *GeocodingClient) GeocodeAddress(adr string) (model.GeocodeResponse, error) {
	params := url.Values{}
	params.Set("address", adr)
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"geocode?"+params.Encode(), nil)
	if err != nil {
		return model.GeocodeResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return model.GeocodeResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.GeocodeResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var geocodeResponse model.GeocodeResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.GeocodeResponse{}, err
	}
	if err := json.Unmarshal(body, &geocodeResponse); err != nil {
		return model.GeocodeResponse{}, err
	}
	return geocodeResponse, nil
}
