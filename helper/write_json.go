package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(writer http.ResponseWriter, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed marshalling json: %w", err)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(bytes)

	if err != nil {
		return fmt.Errorf("failed writing json: %w", err)
	}

	return nil
}
