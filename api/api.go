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

// Series represents a series from the Sonarr API (local sonarr instance API response).
type Series struct {
	ID               int      `json:"id"`
	Title            string   `json:"title"`
	SeasonCount      int      `json:"seasonCount"`
	EpisodesCount    int      `json:"episodesCount"`
	Year             int      `json:"year"`
	Status           string   `json:"status"`
	Path             string   `json:"path"`
	Genres           []string `json:"genres"`
	TotalEpisodes    int      `json:"totalEpisodes"`
	EpisodeFileCount int      `json:"episodeFileCount"`
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

// earches for a specific series by name in the local Sonarr instance
func SonarrLocalSeriesSearch(localSeriesLookupURL, apiKey string, seriesName string) (string, error) {

	/*
		This function performs a lookup for a series that has already been added to the target (local) sonarr instance.

		This API call uses a differnt url endpoint for the search as seen below:

		Example:
		localSeriesLookupURL = "http://<target_instance>:<target_port>/api/series/lookup"

	*/

	// Build the URL to query the series by name
	u, err := url.Parse(localSeriesLookupURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Add the query parameter for the series name
	q := u.Query()
	q.Set("search", seriesName) // Use the query parameter 'search' to search by title
	u.RawQuery = q.Encode()

	// Create a new GET request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set the API key header
	req.Header.Set("X-Api-Key", apiKey)

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

// construct target sonarr instance URL

func ConstructSonarrLocalSeriesURL(sonarrInstance, sonarrPort string) string {
	return fmt.Sprintf("http://%s:%s/api/v3/series", sonarrInstance, sonarrPort)
}

// SonarrFetchAllLocalSeries retrieves all series from the local Sonarr instance.
func SonarrFetchAllLocalSeries(sonarrLocalSeriesUrl, apiKey string) (string, error) {
	// Create a new GET request
	req, err := http.NewRequest("GET", sonarrLocalSeriesUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Add the API key to the headers
	req.Header.Set("X-Api-Key", apiKey)

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

	// Parse the JSON response
	var rawSeries []struct {
		ID         int      `json:"id"`
		Title      string   `json:"title"`
		Status     string   `json:"status"`
		Path       string   `json:"path"`
		Year       int      `json:"year"`
		Genres     []string `json:"genres"`
		Statistics struct {
			SeasonCount       int `json:"seasonCount"`
			TotalEpisodeCount int `json:"totalEpisodeCount"`
			EpisodeFileCount  int `json:"episodeFileCount"`
		} `json:"statistics"`
	}

	err = json.Unmarshal(body, &rawSeries)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Transform the raw data into the Series struct
	var seriesList []Series
	for _, rs := range rawSeries {
		seriesList = append(seriesList, Series{
			ID:            rs.ID,
			Title:         rs.Title,
			Status:        rs.Status,
			Path:          rs.Path,
			Year:          rs.Year,
			Genres:        rs.Genres,
			SeasonCount:   rs.Statistics.SeasonCount,
			EpisodesCount: rs.Statistics.TotalEpisodeCount,
		})
	}

	// Marshal the series list to JSON
	seriesJSON, err := json.Marshal(seriesList)
	if err != nil {
		return "", fmt.Errorf("error marshalling series list to JSON: %w", err)
	}

	return string(seriesJSON), nil
}
