package main

import (
	"encoding/json"
	"fmt"
)

// ClientError defines client error
type ClientError interface {
	Error() string
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

type HTTPError struct {
	Cause  error  `json:"-"`
	Status int    `json:"-"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + ":" + e.Cause.Error()
}

//ResponseBody return response body
func (e *HTTPError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

//ResponseHeaders return status and headers to client
func (e *HTTPError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "problem+application/json;charset=utf-8",
	}
}

//NewHTTPError create new http error
func NewHTTPError(err error, status int, title, detail string) error {
	return &HTTPError{
		Cause:  err,
		Status: status,
		Title:  title,
		Detail: detail,
	}
}
