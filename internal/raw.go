package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetExistingRawEntries[T interface{}](domain string, path string) ([]T, error) {
	url := fmt.Sprintf("https://%s/%s.json", domain, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create raw request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful raw request: %w", err)
	}
	defer resp.Body.Close()

	var response []T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("can't decode raw response: %w", err)
	}

	return response, nil
}
