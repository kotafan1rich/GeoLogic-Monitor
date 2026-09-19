package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return errors.New("invalid request body")
	}
	return nil
}
