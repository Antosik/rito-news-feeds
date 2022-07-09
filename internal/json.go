package internal

import (
	"bytes"
	"encoding/json"
)

func MarshalJSON(data interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(data)

	return buffer.Bytes(), err
}
