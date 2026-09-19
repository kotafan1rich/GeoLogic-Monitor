package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
)

func WriteError(w http.ResponseWriter, err *app.Error) {
	WriteJSON(w, err.Status, err)
}

func WriteServiceError(w http.ResponseWriter, err error) {
	if serviceError, ok := errors.AsType[*app.Error](err); ok {
		WriteError(w, serviceError)
	} else {
		WriteError(w, app.ErrInternal)
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
