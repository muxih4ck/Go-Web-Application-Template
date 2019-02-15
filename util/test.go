package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

// PerformRequest ... Perform Request (without body) util method. Token is optional(pass "")
func PerformRequest(method string, r http.Handler, path string, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// PerformRequestWithBody ... PerformRequest With Body(POST,PUT,DELETE). Token is optional(pass "")
func PerformRequestWithBody(method string, r http.Handler, path string, body []byte, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
