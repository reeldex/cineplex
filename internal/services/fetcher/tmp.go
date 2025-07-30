package fetcher

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type FilterRequest struct {
	ByFormat      string `json:"by_format"`
	ByQuality     string `json:"by_quality"`
	ByAudioFormat string `json:"by_audio_format"`
	ByCinema      string `json:"by_cinema"`
	ByStatus      string `json:"by_status"`
}

func get(cl *http.Client) (CineplexMoviesResponse, error) {
	// Create request payload
	filter := FilterRequest{
		ByFormat:      "",
		ByQuality:     "",
		ByAudioFormat: "",
		ByCinema:      "",
		ByStatus:      "0",
	}

	jsonData, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	// Create request
	req, err := http.NewRequest("POST", "https://cineplex.md/api/getMoviesFiltered", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Set headers (same as above)
	// ... (headers remain the same)

	// Execute request
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	response := CineplexMoviesResponse{}

	return response, json.Unmarshal(respBody, &response)
}
