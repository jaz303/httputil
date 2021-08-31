package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ErrorTypeError       = iota
	ErrorTypeMessageOnly = iota
	ErrorTypeInternal    = iota
)

type ErrorFormatter func(int, int, string, error) interface{}

func defaultErrorFormatter(errorType int, httpStatus int, description string, err error) interface{} {
	message := "an unknown error occurred"
	switch errorType {
	case ErrorTypeError:
		message = fmt.Sprintf("%s (%s)", description, err)
	case ErrorTypeMessageOnly:
		message = description
	case ErrorTypeInternal:
		message = fmt.Sprintf("%s (%s)", description, err)
	}
	return map[string]interface{}{
		"error": message,
	}
}

type JSONResponder struct {
	FormatError ErrorFormatter
}

func NewJSONResponder() *JSONResponder {
	return &JSONResponder{
		FormatError: defaultErrorFormatter,
	}
}

func (r JSONResponder) DecodeJSON(w http.ResponseWriter, b io.Reader, obj interface{}) bool {
	return r.decodeJSON(w, b, obj, false)
}

func (r JSONResponder) DecodeJSONStrict(w http.ResponseWriter, b io.Reader, obj interface{}) bool {
	return r.decodeJSON(w, b, obj, true)
}

func (r JSONResponder) decodeJSON(w http.ResponseWriter, b io.Reader, obj interface{}, strict bool) bool {
	d := json.NewDecoder(b)
	if strict {
		d.DisallowUnknownFields()
	}
	if err := d.Decode(obj); err != nil {
		r.SendError(w, http.StatusBadRequest, "failed to parse request JSON", err)
		return false
	}
	return true
}

func (r JSONResponder) Send(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	return enc.Encode(obj)
}

func (r JSONResponder) SendString(w http.ResponseWriter, status int, json string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(json))
}

func (r JSONResponder) SendEmptyObject(w http.ResponseWriter, status int) {
	r.SendString(w, status, "{}")
}

func (r JSONResponder) SendNotFound(w http.ResponseWriter) {
	r.SendErrorMessage(w, http.StatusNotFound, "not found")
}

func (r JSONResponder) SendUnauthorized(w http.ResponseWriter) {
	r.SendErrorMessage(w, http.StatusUnauthorized, "unauthorized")
}

func (r JSONResponder) SendForbidden(w http.ResponseWriter) {
	r.SendErrorMessage(w, http.StatusForbidden, "forbidden")
}

func (r JSONResponder) SendBadRequestMessage(w http.ResponseWriter, message string) {
	r.SendErrorMessage(w, http.StatusBadRequest, message)
}

func (r JSONResponder) SendError(w http.ResponseWriter, status int, description string, err error) {
	r.Send(w, status, r.FormatError(ErrorTypeError, status, description, err))
}

func (r JSONResponder) SendInternalError(w http.ResponseWriter, description string, err error) {
	r.Send(w, http.StatusInternalServerError, r.FormatError(ErrorTypeInternal, http.StatusInternalServerError, description, err))
}

func (r JSONResponder) SendErrorMessage(w http.ResponseWriter, status int, message string) {
	r.Send(w, status, r.FormatError(ErrorTypeMessageOnly, status, message, nil))
}
