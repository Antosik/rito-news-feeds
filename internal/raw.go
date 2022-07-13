package internal

import (
	"encoding/json"
	"fmt"
)

func GetExistingRawEntries[T interface{}](file FeedFile) ([]T, error) {
	var response []T

	err := json.Unmarshal(file.Buffer, &response)
	if err != nil {
		return nil, fmt.Errorf("can't decode raw response: %w", err)
	}

	return response, nil
}
