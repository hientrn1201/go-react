package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// interface{} means that data can take any kind of type

// helper function to write JSON response and handle any possible error
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	//convert data to JSON
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// add any existing header to the response header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// complete header set up for the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// finally write the JSON to the response
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

// helper function to read JSON from requests and handle any possible error
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 //limit to 1 megabytes
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// create a new decoder that read r.Body
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	// use the decoder to read from r.Body, store JSON value to data pointer variable
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// attempts to decode another JSON value after decode the first JSON to data pointer
	err = dec.Decode(&struct{}{})

	//legit body will return EOF error, otherwise this body contains more than one JSON value
	if err != io.EOF {
		return errors.New("body must only contain one JSON value")
	}

	return nil
}

// helper function to format error when writing and reading JSON
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	// set default status code to 400 - bad request
	statusCode := http.StatusBadRequest

	// if other error code provided, change to that
	if len(status) > 0 {
		statusCode = status[0]
	}

	// create JSON error instance from JSONResponse object
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	// write error in JSON to the response
	return app.writeJSON(w, statusCode, payload)
}
