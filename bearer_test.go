package httputil

import (
	"net/http"
	"testing"
)

func TestBearerParsing(t *testing.T) {
	r := http.Request{}
	r.Header = make(http.Header)
	r.Header.Add("Authorization", "Bearer   asd123  ")

	bt, err := ExtractBearerToken(&r)
	if err != nil || bt != "asd123" {
		t.Errorf("failed to parse bearer token, got `%s`, error: %v", bt, err)
	}
}
