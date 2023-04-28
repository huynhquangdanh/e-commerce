package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type JSONRResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 //1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single json value")
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONRResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func (app *application) genRandomString() string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 1000
	randomNumber := rand.Intn(1000)

	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Concatenate the timestamp and random number
	input := fmt.Sprintf("%d%d", timestamp, randomNumber)

	// Create an MD5 hash of the input string
	hasher := md5.New()
	hasher.Write([]byte(input))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Get the first 8 characters of the hash string
	uniqueString := hash[:8]

	fmt.Println("Generated unique string:", uniqueString)
	return uniqueString
}

func (app *application) GetDiscountRateOnPrice(price int) float64 {
	switch {
	case price > 0 && price <= 5:
		return 0.5
	case price > 5 && price <= 10:
		return 0.6
	case price > 10 && price <= 20:
		return 0.7
	case price > 20:
		return 0.8
	default:
		return 1

	}
}
