package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Struct to match the response from Sonarr's series lookup API
type SeriesLookupResponse struct {
	Title string `json:"title"`
}

// SonarrAPICall performs a GET request to the specified URL with the provided API key, query parameter, and header name.
func SonarrSeriesLookupAPICall(baseURL, apiKeyHeader, apiKey, queryParamName, queryParamValue string) (string, error) {
	// Build the URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Add the query parameter
	q := u.Query()
	q.Set(queryParamName, queryParamValue)
	u.RawQuery = q.Encode()

	// Create a new GET request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set the API key header
	req.Header.Set(apiKeyHeader, apiKey)

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

// Unmarshal and process the API response
func ProcessSeriesLookupResponse(responseBody string) ([]string, error) {
	var seriesResponses []SeriesLookupResponse

	// Unmarshal the JSON response into a slice of structs
	err := json.Unmarshal([]byte(responseBody), &seriesResponses)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling API response: %w", err)
	}

	// Create a slice to hold only the titles
	var titles []string
	for _, series := range seriesResponses {
		titles = append(titles, series.Title)
	}

	return titles, nil
}
