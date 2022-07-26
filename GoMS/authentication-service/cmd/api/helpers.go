package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(res http.ResponseWriter, req *http.Request, data any) error {
	maxBytes := 1048576

	req.Body = http.MaxBytesReader(res, req.Body, int64(maxBytes))

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(res http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			res.Header()[key] = value
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	_, err = res.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(res http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(res, statusCode, payload)
}
