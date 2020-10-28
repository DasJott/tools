package rest

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// Response is the response of a yarc.request
type Response struct {
	// Data is the response data
	Data bytes.Buffer
	// Status is the http status
	Status int
	// Error means local errors
	Error error
}

// Unmarshal pushes the data into a given struct
func (r *Response) Unmarshal(data interface{}) error {
	err := r.Error
	if err == nil {
		err = json.Unmarshal(r.Data.Bytes(), data)
	}
	return err
}

// String returns the response data as a string
func (r *Response) String() string {
	if r.Error == nil {
		return r.Data.String()
	}
	return ""
}

// Int returns the response data as a string
func (r *Response) Int() (n int64, err error) {
	if r.Error == nil {
		n, err = strconv.ParseInt(strings.TrimSpace(r.Data.String()), 10, 64)
	}
	return n, err
}

// Bool returns the response data as a string
func (r *Response) Bool() (b bool, err error) {
	if r.Error == nil {
		b, err = strconv.ParseBool(strings.TrimSpace(r.Data.String()))
	}
	return b, err
}
