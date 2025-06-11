package app

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Error   bool                   `json:"error"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (app *Application) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return err
}

func (app *Application) readJSON(r *http.Request, dataStore interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dataStore)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) errorJSON(w http.ResponseWriter, err error, statusCode int) error {
	var errorJSONData JSONResponse
	errorJSONData.Error = true
	errorJSONData.Message = err.Error()

	return app.writeJSON(w, statusCode, errorJSONData)
}
