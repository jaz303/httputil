package httputil

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeader = errors.New("no auth header")
var ErrBearerFormatError = errors.New("bearer format error")

const bearer = "Bearer "

func ExtractBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", ErrNoAuthHeader
	} else if !strings.HasPrefix(auth, bearer) {
		return "", ErrBearerFormatError
	}
	return strings.TrimSpace(auth[len(bearer):]), nil
}

func ExtractBearerTokenBase64(r *http.Request) ([]byte, error) {
	tok, err := ExtractBearerToken(r)
	if err != nil {
		return nil, err
	}
	bytes, err := base64.StdEncoding.DecodeString(tok)
	if err != nil {
		return nil, ErrBearerFormatError
	}
	return bytes, nil
}
